use reqwest::{blocking::Client, Method};
use serde::de::DeserializeOwned;
use std::io::Write;
use std::*;

#[allow(dead_code)]
pub fn exe_fp() -> path::PathBuf {
    env::current_exe().unwrap()
}

#[allow(dead_code)]
pub fn join_fp<P: AsRef<path::Path>>(paths: &[P]) -> String {
    let mut fp = exe_fp();
    for path in paths {
        fp.push(path);
    }
    fp.to_str().unwrap().to_string()
}

#[allow(dead_code)]
pub fn cover_file(fp: &str, data: &Vec<u8>) {
    let path = path::Path::new(fp);
    if let Some(dir) = path.parent() {
        let _ = fs::create_dir_all(dir);
    }

    fs::write(path, data).unwrap()
}

#[allow(dead_code)]
pub fn cmd_enter(dir: &str, cmdline: &str, shell: &str) -> Result<Vec<u8>, String> {
    let work_dir = if dir.is_empty() {
        env::current_dir().map_err(|e| format!("dir.is_empty/ {}", e))?
    } else {
        path::Path::new(dir).to_path_buf()
    };

    let mut cmd = process::Command::new(shell)
        .current_dir(work_dir)
        .stdin(process::Stdio::piped())
        .stdout(process::Stdio::piped())
        .stderr(process::Stdio::piped())
        .spawn()
        .map_err(|e| e.to_string())?;

    {
        let stdin = cmd.stdin.as_mut().unwrap();
        stdin
            .write_all(cmdline.as_bytes())
            .map_err(|e| e.to_string())?;
        stdin.write_all(b"\nexit\n").map_err(|e| e.to_string())?;
    }

    let output = cmd.wait_with_output().map_err(|e| e.to_string())?;

    if !output.status.success() {
        eprintln!("stderr: command failed: {}", output.status);
    }

    Ok(output.stdout)
}

#[allow(dead_code)]
pub fn cat_string(fp: &str) -> String {
    let b = cat_bytes(fp);
    ts(&b)
}

#[allow(dead_code)]
pub fn cat_bytes(fp: &str) -> Vec<u8> {
    let b = fs::read(fp);
    if b.is_err() {
        stderr(b.err().unwrap());
        return b"".to_vec();
    }
    b.unwrap()
}

#[allow(dead_code)]
pub fn file_exist(fp: &str) -> bool {
    path::Path::new(fp).exists()
}

#[allow(dead_code)]
pub fn ts(a: &Vec<u8>) -> String {
    unsafe { str::from_utf8_unchecked(a.as_slice()).to_string() }
}

fn stderr(err: impl fmt::Display) {
    eprintln!("stderr: {}", err);
}

#[allow(dead_code)]
pub fn unmarshal<T: DeserializeOwned>(b: &[u8], v: &mut T) -> Result<(), String> {
    *v = serde_json::from_slice(b).map_err(|e| e.to_string())?;
    Ok(())
}

#[macro_export]
macro_rules! pc {
    () => {{
        fn f() {}
        fn type_name_of<T>(_: T) -> &'static str {
            any::type_name::<T>()
        }
        let name = type_name_of(f);
        name.strip_suffix("::f").unwrap_or(name)
    }};
}

#[macro_export]
macro_rules! to_json_bytes {
    ($v:expr) => {
        serde_json::to_vec($v).unwrap()
    };
    ($v:expr, $format:expr) => {
        if $format {
            serde_json::to_vec_pretty($v).unwrap()
        } else {
            serde_json::to_vec($v).unwrap()
        }
    };
}

#[allow(dead_code)]
pub fn http_do(method: Method, url: &str) -> Result<(reqwest::header::HeaderMap, Vec<u8>), String> {
    // 创建忽略证书验证的阻塞客户端
    let client = Client::builder()
        .danger_accept_invalid_certs(true) // 忽略证书验证
        .danger_accept_invalid_hostnames(true) // 忽略主机名验证
        .build()
        .map_err(|e| e.to_string())?;

    // 发送阻塞请求
    let response = client
        .request(method, url)
        .send()
        .map_err(|e| e.to_string())?;

    // 获取响应状态码
    let status = response.status();

    // 处理响应
    let headers = response.headers().clone();
    let body = response.bytes().map_err(|e| e.to_string())?;
    let body = body.to_vec();

    // 检查HTTP状态码
    if !status.is_success() {
        return Err(format!(
            "HTTP error: {} - {}",
            status,
            String::from_utf8_lossy(&body)
        )
        .into());
    }

    Ok((headers, body))
}

#[allow(dead_code)]
pub fn unzip(fp: &str, out: &str) -> Result<(), String> {
    let source = fs::File::open(fp).map_err(|err| format!("fs::File::open {} {}", fp, err))?;

    let mut archive = zip::read::ZipArchive::new(source)
        .map_err(|err| format!("zip::read::ZipArchive::new {} {}", fp, err))?;

    for i in 0..archive.len() {
        let file = archive.by_index(i);
        if file.is_err() {
            println!("archive.by_index {} {}", i, file.err().unwrap());
            continue;
        }
        let mut file = file.unwrap();

        let out_fp = path::Path::new(out).join(file.mangled_name());

        if file.is_dir() {
            let err = fs::create_dir_all(&out_fp);
            if err.is_err() {
                println!("fs::create_dir_all {} {}", i, err.err().unwrap());
                continue;
            }
        } else {
            if let Some(p) = out_fp.parent() {
                if !p.exists() {
                    let err = fs::create_dir_all(p);
                    if err.is_err() {
                        println!("fs::create_dir_all {} {}", i, err.err().unwrap());
                        continue;
                    }
                }
            }

            let out_file = fs::File::create(&out_fp);
            if out_file.is_err() {
                println!(
                    "fs::File::create {} {}",
                    out_fp.to_str().unwrap(),
                    out_file.err().unwrap()
                );
                continue;
            }
            let mut out_file = out_file.unwrap();

            let err = io::copy(&mut file, &mut out_file);
            if err.is_err() {
                println!(
                    "io::copy {} {}",
                    out_fp.to_str().unwrap(),
                    err.err().unwrap()
                );
                continue;
            }
        }
    }

    Ok(())
}
