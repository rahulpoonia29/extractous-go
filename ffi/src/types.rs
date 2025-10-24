use std::os::raw::{c_char, c_int};

#[repr(C)]
pub struct CExtractor {
    _private: [u8; 0],
    _marker: core::marker::PhantomData<(*mut u8, core::marker::PhantomPinned)>,
}
#[repr(C)]
pub struct CStreamReader {
    _private: [u8; 0],
    _marker: core::marker::PhantomData<(*mut u8, core::marker::PhantomPinned)>,
}
#[repr(C)]
pub struct CPdfParserConfig {
    _private: [u8; 0],
    _marker: core::marker::PhantomData<(*mut u8, core::marker::PhantomPinned)>,
}
#[repr(C)]
pub struct COfficeParserConfig {
    _private: [u8; 0],
    _marker: core::marker::PhantomData<(*mut u8, core::marker::PhantomPinned)>,
}
#[repr(C)]
pub struct CTesseractOcrConfig {
    _private: [u8; 0],
    _marker: core::marker::PhantomData<(*mut u8, core::marker::PhantomPinned)>,
}

#[repr(C)]
pub struct CMetadata {
    /// Array of pointers to null-terminated key strings
    pub keys: *mut *mut c_char,
    /// Array of pointers to null-terminated value strings
    pub values: *mut *mut c_char,
    /// The number of key-value pairs in the arrays
    pub len: libc::size_t,
}

pub const CHARSET_UTF_8: c_int = 0;
pub const CHARSET_US_ASCII: c_int = 1;
pub const CHARSET_UTF_16BE: c_int = 3;

pub const PDF_OCR_STRATEGY_NO_OCR: c_int = 0;
pub const PDF_OCR_STRATEGY_OCR_ONLY: c_int = 1;
pub const PDF_OCR_STRATEGY_OCR_AND_TEXT_EXTRACTION: c_int = 2;
pub const PDF_OCR_STRATEGY_AUTO: c_int = 3;
