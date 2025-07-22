#!/usr/bin/env -S uv run
# /// script
# dependencies = [
#     "httpx==0.27.2",
#     "click==8.1.7",
#     "rich==13.9.4",
# ]
# ///

from typing import Optional, Tuple
import httpx
import click
from rich.progress import Progress, BarColumn, DownloadColumn, TimeRemainingColumn, TaskID
from pathlib import Path


def format_bytes(size: int) -> str:
    """Convert bytes to human readable format."""
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if size < 1024.0:
            return f"{size:.1f} {unit}"
        size /= 1024.0
    return f"{size:.1f} PB"


@click.command()
@click.argument('url', type=str)
@click.option('--output', '-o', type=click.Path(), help='Output filename (default: inferred from URL)')
@click.option('--chunk-size', type=int, default=8192, help='Download chunk size in bytes')
def download(url: str, output: Optional[str], chunk_size: int) -> None:
    """Download a file from URL with a progress bar."""
    
    # Determine output filename
    output_path: Path
    if not output:
        filename: str = url.split('/')[-1]
        if not filename or filename == '':
            filename = 'download'
        output_path = Path(filename)
    else:
        output_path = Path(output)
    
    # Create progress bar
    with Progress(
        "[progress.description]{task.description}",
        BarColumn(),
        DownloadColumn(),
        TimeRemainingColumn(),
    ) as progress:
        
        with httpx.stream("GET", url, follow_redirects=True) as response:
            response.raise_for_status()
            
            # Get total file size
            content_length: Optional[str] = response.headers.get('content-length')
            total_size: int = int(content_length) if content_length else 0
            
            # Create progress task
            task: TaskID = progress.add_task(f"Downloading {output_path.name}", total=total_size)
            
            # Download and write file
            with open(output_path, 'wb') as f:
                chunk: bytes
                for chunk in response.iter_bytes(chunk_size):
                    f.write(chunk)
                    progress.update(task, advance=len(chunk))
    
    file_size: int = output_path.stat().st_size
    click.echo(f"\n✓ Downloaded to {output_path}")
    click.echo(f"  Size: {format_bytes(file_size)}")


if __name__ == '__main__':
    download()