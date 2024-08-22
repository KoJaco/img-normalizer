# Image Normalizer CLI Tool

## Overview

The **Image Normalizer CLI Tool** is a command-line application built in Go that processes images to adjust them to the nearest common aspect ratio. The tool can handle individual image files or directories containing images, including nested directories. It preserves the directory structure when saving the processed images and generates a comprehensive log of all processed images.

## Features

-   Adjusts images to the nearest common aspect ratio.
-   Supports JPG, PNG, and WebP image formats.
-   Allows specifying a tolerance for how much an image can be scaled or cropped.
-   Handles single image files and directories (including nested directories).
-   Preserves the directory structure when saving processed images.
-   Generates a CSV log file recording details of all processed images, including original dimensions, the chosen aspect ratio, new dimensions, and any issues encountered.

## Nearest Common Aspect Ratios:

-   1:1
-   4:3
-   4:5
-   3:2
-   16:9
-   21:9
-   9:16

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/img-normalizer.git
    cd img-normalizer
    ```

2. **Install dependencies:**

    Ensure you have Go installed. If not, download it from [golang.org](https://golang.org/dl/).

    ```bash
    go mod tidy
    ```

3. **Build the CLI tool:**

    ```bash
    go build -o img-normalizer ./cmd/img-normalizer
    ```

## Usage

### Command-line Options

```bash
./img-normalizer --src=<source_path> --dest=<destination_path> --tolerance=<pixel_tolerance>
```

-   `--src` (required): The source path of either a directory or a single image file. If the path is a single image file, make sure to include the extension (.jpg, .png, or .webp)

-   `--dest` (required): The destination directory where processed images and a log of the whole process will be saved to.

-   `--tolerance` (optional): The maximum number of pixels that can be adjusted for either the width or the height without changing the image dimensions. The default is `100`.

## Example Commands

### Process a Single Image:

```bash
./img-normalizer --src="example/images/raw/sample.jpg" --dest="example/images/out" --tolerance=100
```

### Process all images in a directory (Flat Structure) with a tolerance of 50:

```bash
./img-normalizer --src="example/images/raw" --dest="example/images/out" --tolerance=50
```

### Process all images in a directory (Nested Structure):

```bash
./img-normalizer --src="example/images/nested_raw" --dest="example/images/out" --tolerance=100
```

## Log File

The tool generates an `image_process_log.csv` file in the destination directory. This log includes:

-   Image Path: The relative or absolute path of the image.
-   Original Dimensions: The original width and height of the image.
-   Chosen Aspect Ratio: The aspect ratio (e.g., "3:2") chosen as the closest common aspect of the image.
-   New Dimensions: The new width and height after adjustment.
-   Status: indicates whether the image was adjusted, not adjusted (with a reason), or if there was an error.

## Supported Image Formats

-   JPG/JPEG (.jpg or .jpeg)
-   PNG (.png)
-   WebP (.webp)

## Timeline

1. Add flag to target specific file type(s) in the supported image format list
2. Add support for additional image formats
3. Allow users to append their own list of common aspect ratios (maybe)
4. Make a web service where you can just drag files/folders in and get back a zip download for processed images (maybe)

## Contributing

If you have ideas for improvements or find bugs, feel free to open an issue or submit a pull request.

## license

This project is licensed under the MIT license. See the `LICENSE` file for details.
