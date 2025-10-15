use std::env;
use std::fs;
use std::path::PathBuf;

fn main() {
    // Skip build during docs.rs documentation builds
    if env::var("DOCS_RS").is_ok() {
        return;
    }

    let manifest_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let root_dir = PathBuf::from(&manifest_dir);

    println!("cargo:warning=Building extractous-ffi");

    // 1. Generate C header file using cbindgen
    generate_header(&manifest_dir, &root_dir);

    // 2. Set RPATH for runtime library discovery (platform-specific)
    set_rpath();

    // 3. Configure rerun triggers for build script
    configure_rerun_triggers();
}

/// Generate C header file from Rust source using cbindgen
///
/// This function reads the cbindgen.toml configuration and generates
/// a C header file that can be used with Go cgo or any C-compatible FFI.
fn generate_header(crate_dir: &str, root_dir: &PathBuf) {
    let include_dir = root_dir.join("include");

    // Create include directory if it doesn't exist
    fs::create_dir_all(&include_dir).expect("Failed to create include directory");

    let header_path = include_dir.join("extractous.h");

    // Use cbindgen with the configuration file
    match cbindgen::Builder::new()
        .with_crate(crate_dir)
        .with_config(
            cbindgen::Config::from_file("cbindgen.toml").unwrap_or_else(|_| {
                println!("cargo:warning=cbindgen.toml not found, using default configuration");
                cbindgen::Config::default()
            }),
        )
        .generate()
    {
        Ok(bindings) => {
            bindings.write_to_file(&header_path);
            println!(
                "cargo:warning=Successfully generated C header: {}",
                header_path.display()
            );
            println!(
                "cargo:warning=Header file size: {} bytes",
                fs::metadata(&header_path).map(|m| m.len()).unwrap_or(0)
            );
        }
        Err(e) => {
            println!("cargo:warning=Failed to generate C header: {:?}", e);
            println!("cargo:warning=Continuing build without header generation");
        }
    }
}

/// Set RPATH for runtime library discovery
///
/// Configures the dynamic linker to find shared libraries at runtime.
/// This is platform-specific and ensures that the Go bindings can find
/// the extractous FFI library and any dependencies.
fn set_rpath() {
    let target = env::var("TARGET").unwrap();

    println!("cargo:warning=Configuring RPATH for target: {}", target);

    if target.contains("linux") {
        // Linux: Use $ORIGIN to reference the library's directory
        // This allows the library to be relocatable
        println!("cargo:rustc-link-arg=-Wl,-rpath,$ORIGIN");
        println!("cargo:rustc-link-arg=-Wl,-z,origin");
        // Disable new dtags to use RPATH instead of RUNPATH
        // This makes the path non-overridable by LD_LIBRARY_PATH
        println!("cargo:rustc-link-arg=-Wl,--disable-new-dtags");
        println!("cargo:warning=Linux RPATH configured with $ORIGIN");
    } else if target.contains("darwin") || target.contains("macos") {
        // macOS: Use @loader_path to reference the library's directory
        println!("cargo:rustc-link-arg=-Wl,-rpath,@loader_path");
        // Set install name to use @rpath for better relocatability
        println!("cargo:rustc-link-arg=-Wl,-install_name,@rpath/libextractous_ffi.dylib");
        println!("cargo:warning=macOS RPATH configured with @loader_path");
    } else if target.contains("windows") {
        // Windows: DLLs are searched in the current directory by default
        // No special RPATH configuration needed
        println!("cargo:warning=Windows target detected - using default DLL search path");
    } else {
        println!("cargo:warning=Unknown target platform, skipping RPATH configuration");
    }
}

/// Configure build script rerun triggers
///
/// Tells Cargo when to rerun this build script.
fn configure_rerun_triggers() {
    // Rerun if source files change
    println!("cargo:rerun-if-changed=src");

    // Rerun if build script changes
    println!("cargo:rerun-if-changed=build.rs");

    // Rerun if cbindgen config changes
    println!("cargo:rerun-if-changed=cbindgen.toml");

    // Rerun if Cargo.toml changes (dependencies, features, etc.)
    println!("cargo:rerun-if-changed=Cargo.toml");
}
