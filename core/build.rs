use std::env;
use std::path::PathBuf;

fn main() {
    // Skip when building docs
    if env::var("DOCS_RS").is_ok() {
        return;
    }

    println!("cargo:warning=Using pre-built libtika_native libraries");

    // Get paths
    let manifest_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let root_dir = PathBuf::from(&manifest_dir).parent().unwrap().to_path_buf();

    // Detect platform
    let target_os = env::var("CARGO_CFG_TARGET_OS").unwrap();
    let target_arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();

    let platform_dir = match (target_os.as_str(), target_arch.as_str()) {
        ("linux", "x86_64") => "linux_amd64",
        ("linux", "aarch64") => "linux_arm64",
        ("macos", "x86_64") => "darwin_amd64",
        ("macos", "aarch64") => "darwin_arm64",
        ("windows", "x86_64") => "windows_amd64",
        _ => panic!("Unsupported platform: {} {}", target_os, target_arch),
    };

    let libs_dir = root_dir.join("native").join(platform_dir);

    println!("cargo:rustc-link-search=native={}", libs_dir.display());

    // Link to the library
    let lib_name = if target_os == "windows" {
        "tika_native"
    } else {
        "tika_native"
    };
    println!("cargo:rustc-link-lib=dylib={}", lib_name);

    // Verify libraries exist
    if !libs_dir.exists() {
        panic!(
            "\n\n================================================================\n\
             ERROR: Pre-built native libraries not found!\n\
             ================================================================\n\
             Expected at: {}\n\n\
             Please extract libraries from Python wheel:\n\
             1. pip download extractous==0.2.1 --platform manylinux_2_31_x86_64 --only-binary=:all: -d /tmp\n\
             2. unzip /tmp/extractous-*.whl -d /tmp/wheel\n\
             3. mkdir -p {}\n\
             4. cp /tmp/wheel/extractous/*.so {}\n\
             ================================================================\n",
            libs_dir.display(),
            libs_dir.display(),
            libs_dir.display()
        );
    }

    println!("cargo:rustc-link-search=native={}", libs_dir.display());
    println!("cargo:rustc-link-lib=dylib=tika_native");

    // Set rpath for runtime library loading
    if target_os == "linux" {
        println!("cargo:rustc-link-arg=-Wl,-rpath,$ORIGIN");
    } else if target_os == "macos" {
        println!("cargo:rustc-link-arg=-Wl,-rpath,@loader_path");
    }

    // Rerun if libraries change
    println!("cargo:rerun-if-changed={}", libs_dir.display());
}
