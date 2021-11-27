use std::env;
use std::path::Path;
use std::path::PathBuf;
use structopt::StructOpt;

#[derive(StructOpt, Debug)]
#[structopt(name = "basic")]
struct Opt {
    /// Activate ASCII escapes mode
    #[structopt(short, long)]
    escape: bool,

    /// Files to print dir from
    #[structopt(name = "FILE")]
    file: Option<String>,
}

fn main() {
    let args = Opt::from_args();
    let file = args.file;
    let mut path = env::current_dir().expect("CWD is not set");
    if file.is_some() {
        let file_path = file.unwrap();
        check_exsits(&file_path);
        path.push(handle_dot_for_current_dir(&file_path));
    };
    if env::consts::OS == "windows" && args.escape {
        println!("{}", windows_path_escape(&path.to_str().unwrap_or("")));
    } else {
        println!("{}", path.display());
    }
}

fn handle_dot_for_current_dir(n: &String) -> PathBuf {
    match &n[..2] {
        "./" => Path::new(&n.replace("./", "")).to_owned(),
        r#".\"# => Path::new(&n.replace(r#".\"#, "")).to_owned(),
        _ => Path::new(n).to_owned(),
    }
}

fn check_exsits(path: &String) -> () {
    let file_path = Path::new(path);
    if !file_path.exists() {
        panic!("File or DIR {} doesn't exists", file_path.display());
    }
}

fn windows_path_escape(s: &str) -> String {
    s.replace(r#"\"#, r#"\\"#)
}
