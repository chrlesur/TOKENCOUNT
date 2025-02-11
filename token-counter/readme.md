# tokencount

A simple tool to count tokens in files using a Go native library.

## Installation

1.  Install Go: <https://go.dev/dl/>
2.  Clone the repository:

    ```bash
    git clone [repository URL]
    ```
    (Replace `[repository URL]` with the actual repository URL.)
3.  Navigate to the project directory:

    ```bash
    cd token-counter
    ```

## Compilation

```bash
go build -o token-count
```

This will create an executable file named `token-count` in the project directory.

## Usage

```
./token-count [flags] [files...]
```

## Flags

*   `--version`: Displays the version of the application.
*   `--help`: Displays help information.
*   `--debug`: Enables debug mode.
*   `--log-file`: Specifies the log file.
*   `--threads`: Specifies the number of threads to use.
*   `--recursive`: Enables recursive directory processing.

## Example

```
./token-count --recursive --threads=4 --log-file=test.log testdir test3.txt
```

## License

This software is licensed under the GNU General Public License Version 3.
