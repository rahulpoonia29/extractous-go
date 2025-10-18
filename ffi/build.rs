use std::env;
use std::fs;
use std::path::PathBuf;

fn main() {
    // Skip during docs builds
    if env::var("DOCS_RS").is_ok() {
        return;
    }

    let manifest_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let target = env::var("TARGET").unwrap();
    let profile = env::var("PROFILE").unwrap();
    
    println!("cargo:warning=Building extractous-ffi for target: {}", target);
    println!("cargo:warning=Profile: {}", profile);

    // 1. Generate C header
    generate_header(&manifest_dir);

    // 2. Configure RPATH for runtime library discovery
    configure_rpath(&target);

    // 3. Ensure extractous dependency built libraries are discoverable
    setup_extractous_libs(&target, &profile);

    // 4. Configure rerun triggers
    configure_rerun_triggers();
}

fn generate_header(crate_dir: &str) {
    let root_dir = PathBuf::from(crate_dir).parent().unwrap().to_path_buf();
    let include_dir = root_dir.join("include");
    fs::create_dir_all(&include_dir).expect("Failed to create include directory");

    let header_path = include_dir.join("extractous.h");

    match cbindgen::Builder::new()
        .with_crate(crate_dir)
        .with_config(
            cbindgen::Config::from_file("cbindgen.toml")
                .unwrap_or_else(|_| cbindgen::Config::default()),
        )
        .generate()
    {
        Ok(bindings) => {
            bindings.write_to_file(&header_path);
            println!("cargo:warning=Generated C header: {}", header_path.display());
        }
        Err(e) => {
            println!("cargo:warning=Failed to generate header: {:?}", e);
        }
    }
}

fn configure_rpath(target: &str) {
    if target.contains("linux") {
        // Use $ORIGIN for relocatable libraries
        println!("cargo:rustc-link-arg=-Wl,-rpath,$ORIGIN");
        println!("cargo:rustc-link-arg=-Wl,-z,origin");
        println!("cargo:rustc-link-arg=-Wl,--disable-new-dtags");
        println!("cargo:warning=Configured Linux RPATH with $ORIGIN");
    } else if target.contains("darwin") || target.contains("macos") {
        // Use @loader_path for macOS
        println!("cargo:rustc-link-arg=-Wl,-rpath,@loader_path");
        println!("cargo:rustc-link-arg=-Wl,-install_name,@rpath/libextractous_ffi.dylib");
        println!("cargo:warning=Configured macOS RPATH with @loader_path");
    } else if target.contains("windows") {
        println!("cargo:warning=Windows: Using default DLL search path");
    }
}

fn setup_extractous_libs(target: &str, profile: &str) {
    // The extractous crate builds libtika_native via its build.rs
    // We need to ensure those libraries are found during linking
    
    let out_dir = env::var("OUT_DIR").unwrap();
    let target_dir = PathBuf::from(&out_dir)
        .parent().unwrap()
        .parent().unwrap()
        .parent().unwrap()
        .to_path_buf();
    
    // Search for extractous build output
    let build_dir = target_dir.join("build");
    
    if let Ok(entries) = fs::read_dir(&build_dir) {
        for entry in entries.flatten() {
            let path = entry.path();
            if let Some(name) = path.file_name() {
                if name.to_str().unwrap().starts_with("extractous-") {
                    let libs_dir = path.join("out").join("libs");
                    if libs_dir.exists() {
                        println!("cargo:rustc-link-search={}", libs_dir.display());
                        println!("cargo:warning=Found extractous libs: {}", libs_dir.display());
                    }
                }
            }
        }
    }
}

fn configure_rerun_triggers() {
    println!("cargo:rerun-if-changed=src");
    println!("cargo:rerun-if-changed=build.rs");
    println!("cargo:rerun-if-changed=cbindgen.toml");
    println!("cargo:rerun-if-changed=Cargo.toml");
}
