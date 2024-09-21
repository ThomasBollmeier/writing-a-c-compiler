use clap::{Args, Parser};
use regex::Regex;
use std::io::Result;
use std::process::Command;

#[derive(Parser)]
#[command(
    author = clap::crate_authors!("\n"),
    version = clap::crate_version!(),
    about = "TBCC - Thomas Bollmeier's C Compiler"
)]
struct Options {
    #[arg(help = "The input file to compile")]
    input_file: String,
    #[command(flatten)]
    mode: CompilerMode,
}

#[derive(Args, Clone)]
#[group(multiple = false)]
struct CompilerMode {
    #[arg(long, help = "Run Lexer on the input file and exit")]
    lex: bool,
    #[arg(long, help = "Parse the input file and exit")]
    parse: bool,
    #[arg(
        long,
        help = "Lex, parse and assemble the input file and exit before code emission"
    )]
    codegen: bool,
    #[arg(short = 'S', help = "Compile the input file and exit")]
    create_assembly: bool,
}

fn main() -> Result<()> {
    let options = Options::parse();

    let preprocessed_file = preprocess_file(&options.input_file)?;

    let assembly_file = compile_file(&preprocessed_file)?;
    remove_file(&preprocessed_file)?;

    let exec_file = assemble_link_file(&assembly_file)?;
    remove_file(&assembly_file)?;

    println!("Compiled to {}", exec_file);

    Ok(())
}

fn strip_extension(file_name: &str) -> String {
    let re = Regex::new(r"\.[^.]+$").unwrap();
    re.replace(file_name, "").to_string()
}

fn preprocess_file(source_file: &str) -> Result<String> {
    let base_name = strip_extension(source_file);
    let output_file = format!("{}.i", base_name);
    let output = std::process::Command::new("gcc")
        .arg("-E")
        .arg("-P")
        .arg(source_file)
        .arg("-o")
        .arg(&output_file)
        .output()?;

    if !output.status.success() {
        panic!("Failed to preprocess file");
    }

    Ok(output_file)
}

fn compile_file(preprocessed_file: &str) -> Result<String> {
    let base_name = strip_extension(preprocessed_file);
    let output_file = format!("{}.s", base_name);
    let output = Command::new("gcc")
        .arg("-S")
        .arg(preprocessed_file)
        .arg("-o")
        .arg(&output_file)
        .output()?;

    if !output.status.success() {
        panic!("Failed to compile file");
    }

    Ok(output_file)
}

fn assemble_link_file(assembly_file: &str) -> Result<String> {
    let exec_file = strip_extension(assembly_file);
    let output = Command::new("gcc")
        .arg(assembly_file)
        .arg("-o")
        .arg(&exec_file)
        .output()?;

    if !output.status.success() {
        panic!("Failed to assemble file");
    }

    Ok(exec_file)
}

fn remove_file(file_name: &str) -> Result<()> {
    std::fs::remove_file(file_name)
}
