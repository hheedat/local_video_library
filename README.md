# Local Video Library

A web-based video library that allows you to browse and play local video files with their corresponding poster images.

## Features

- Scan local folders for video files and their corresponding poster images
- YouTube-like interface for browsing videos
- Search functionality to find videos by name
- Responsive design that works on both desktop and mobile devices

## Requirements

- Python 3
- Flask

## Installation

1. Clone this repository or download the files
2. Install the required packages:
   ```bash
   pip install -r requirements.txt
   ```

## Usage

1. Run the application:
   ```bash
   python app.py
   ```
2. Add folder list in `dir.conf`, and the application will scan the folder for video files (like .mp4) and their corresponding poster images (like .jpg)

3. Open your web browser and navigate to `http://localhost:5000`

4. Use the search bar to find specific videos by name

5. Click on any video card to play the video in a modal window

## File Structure

The application expects video files and their corresponding poster images to follow this naming convention (just need has same prefix):
- Video file: `VIDEO_NAME.mp4`
- Poster image: `VIDEO_NAME.jpg`

For example:
- `XXX-000.mp4`
- `XXX-000.jpg`

## Notes
- The application runs locally and does not upload any files to external servers
