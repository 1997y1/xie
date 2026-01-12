use file_rotate::{compression::Compression, suffix::AppendCount};
use std::io::Write;
use std::*;

#[derive(Copy, Clone)]
struct Options {
    dr: bool,
    d: bool,
    s: bool,
    i: bool,
    w: bool,
    e: bool,
    print: bool,
}

static CONTROL: Options = Options {
    dr: true,
    d: true,
    s: true,
    i: true,
    w: true,
    e: true,
    print: true,
};

#[allow(dead_code)]
pub struct RotatingLogger {
    path: path::PathBuf,
    writer: file_rotate::FileRotate<AppendCount>,
}

impl RotatingLogger {
    pub fn new(
        path: impl Into<path::PathBuf>,
        max_size: u64,
        max_backups: usize,
    ) -> io::Result<Self> {
        let path: path::PathBuf = path.into();
        let writer = file_rotate::FileRotate::new(
            path.clone(),
            AppendCount::new(max_backups),
            file_rotate::ContentLimit::BytesSurpassed(max_size as usize),
            Compression::None,
            None,
        );
        Ok(Self { path, writer })
    }
    fn write(&mut self, buf: &[u8]) -> io::Result<()> {
        self.writer.write_all(buf)
    }
}

pub static LOGGER: sync::OnceLock<sync::Mutex<RotatingLogger>> = sync::OnceLock::new();

pub fn error(txt: &str, pc: &str) {
    if CONTROL.e {
        log_kernel("ERROR  ", txt.as_bytes(), CONTROL.print, pc);
    }
}

pub fn debug(txt: &str, pc: &str) {
    if CONTROL.d {
        log_kernel("DEBUG  ", txt.as_bytes(), CONTROL.print, pc);
    }
}

pub fn warn(txt: &str, pc: &str) {
    if CONTROL.w {
        log_kernel("WARN   ", txt.as_bytes(), CONTROL.print, pc);
    }
}

pub fn info(txt: &str, pc: &str) {
    if CONTROL.i {
        log_kernel("INFO   ", txt.as_bytes(), CONTROL.print, pc);
    }
}

pub fn state(txt: &str) {
    if CONTROL.s {
        log_kernel("STATE  ", txt.as_bytes(), CONTROL.print, "");
    }
}

pub fn raw_debug(txt: &str) {
    if CONTROL.dr {
        println!("{}", txt);
        let txt = txt.to_string() + "\n";
        let logger = LOGGER.get();
        if !logger.is_none() {
            let _ = logger.unwrap().lock().unwrap().write(&txt.as_bytes());
        }
    }
}

#[track_caller]
fn log_kernel(tag: &str, data: &[u8], output: bool, func: &str) {
    let ts = chrono::Local::now().format("%m%d %H:%M:%S%.3f").to_string();
    let pc = panic::Location::caller();
    let file = pc.file();
    let line = pc.line();

    // Get code file.
    let file = file.rsplit('/').next().unwrap_or(file);

    // Get log prefix.
    let prefix = if !func.is_empty() {
        format!("{}{} ({}) {}:{} ", tag, ts, func, file, line)
    } else {
        format!("STATE  {} ", ts)
    };

    // Write buffer.
    let mut buf = Vec::with_capacity(1024);
    buf.extend_from_slice(prefix.as_bytes());
    buf.extend_from_slice(data);
    buf.push(b'\n');

    // Print to console.
    if output {
        let _ = io::stdout().write_all(&buf);
    }

    let logger = LOGGER.get();
    if !logger.is_none() {
        let _ = logger.unwrap().lock().unwrap().write(&buf);
    }
}

// use std::*;
//
// fn main() {
//     let logger =
//         console::RotatingLogger::new("console.log", 100 * 1024 * 1024, 10).expect("open console.log");
//     let _ = console::LOGGER.set(sync::Mutex::new(logger));
//
//     kei_my_boy();
// }
//
// fn kei_my_boy() {
//     let txt = "error message.";
//
//     console::error(txt,xie::pc!());
//     console::debug(txt,xie::pc!());
//     console::warn(txt,xie::pc!());
//     console::info(txt,xie::pc!());
//     console::state(txt);
//     console::raw_debug(txt);
// }

