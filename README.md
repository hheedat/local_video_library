# Local Video Library

A web-based video library that allows you to browse and play local video files with their corresponding poster images.

## Features

- Scan local folders for video files and their corresponding poster images
- YouTube-like interface for browsing videos
- Search functionality to find videos by name
- Real-time file system monitoring for automatic updates
- Responsive design that works on both desktop and mobile devices

## Requirements

- Python 3.7 or higher
- Flask
- watchdog

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

2. Open your web browser and navigate to `http://localhost:5000`

3. Click the "Add Folder" button to add a folder containing your videos and poster images

4. The application will scan the folder for video files (.mp4) and their corresponding poster images (.jpg)

5. Use the search bar to find specific videos by name

6. Click on any video card to play the video in a modal window

## File Structure

The application expects video files and their corresponding poster images to follow this naming convention:
- Video file: `VIDEO_NAME.mp4`
- Poster image: `VIDEO_NAME.jpg`

For example:
- `XXX-000.mp4`
- `XXX-000.jpg`

## Notes

- The application monitors the added folders for changes, so any new videos or poster images added to the folders will automatically appear in the interface
- Video files must be in MP4 format
- Poster images must be in JPG format
- The application runs locally and does not upload any files to external servers
