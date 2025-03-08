import os
import json
from flask import Flask, render_template, jsonify, request, send_file
from pathlib import Path
import mimetypes
import re

app = Flask(__name__)

# Global variables to store video data
video_data = {}

# Configure video directories
VIDEO_DIRECTORIES = [
    r'F:\xxxxx',
    r'E:\xxxxx'
]

# Supported file extensions
VIDEO_EXTENSIONS = {'.mp4', '.mkv', '.avi', '.mov', '.wmv', '.flv', '.webm'}
IMAGE_EXTENSIONS = {'.jpg', '.jpeg', '.png', '.gif', '.webp'}

def get_base_name(filename):
    """Extract base name from filename, handling various suffixes."""
    # Remove file extension
    name = Path(filename).stem
    # Remove common suffixes
    suffixes = [r'\.FHD$', r'\.HD$', r'\.FHD\.', r'\.HD\.']
    for suffix in suffixes:
        name = re.sub(suffix, '', name)
    return name

def scan_directory(directory):
    """Scan a directory for video and poster files."""
    directory = Path(directory)
    if not directory.exists():
        print(f"Directory not found: {directory}")
        return
    
    print(f"\nScanning directory: {directory}")
    print("-" * 50)
    
    for file_path in directory.rglob('*'):
        if file_path.is_file():
            suffix = file_path.suffix.lower()
            if suffix in VIDEO_EXTENSIONS or suffix in IMAGE_EXTENSIONS:
                base_name = get_base_name(file_path.name)
                if base_name not in video_data:
                    video_data[base_name] = {'video': None, 'poster': None}
                
                if suffix in VIDEO_EXTENSIONS:
                    video_data[base_name]['video'] = str(file_path)
                    print(f"Found video: {file_path.name} (base: {base_name})")
                elif suffix in IMAGE_EXTENSIONS:
                    video_data[base_name]['poster'] = str(file_path)
                    print(f"Found poster: {file_path.name} (base: {base_name})")
    
    print("-" * 50)

def print_scan_summary():
    """Print summary of scanned files."""
    total_videos = sum(1 for data in video_data.values() if data['video'])
    total_posters = sum(1 for data in video_data.values() if data['poster'])
    complete_pairs = sum(1 for data in video_data.values() if data['video'] and data['poster'])
    
    print("\nScan Summary:")
    print("-" * 50)
    print(f"Total videos found: {total_videos}")
    print(f"Total posters found: {total_posters}")
    print(f"Complete pairs (video + poster): {complete_pairs}")
    
    # Check for mismatches
    videos_without_posters = [name for name, data in video_data.items() if data['video'] and not data['poster']]
    posters_without_videos = [name for name, data in video_data.items() if not data['video'] and data['poster']]
    
    if videos_without_posters:
        print("\nVideos without posters:")
        for name in videos_without_posters:
            print(f"- {name}")
    
    if posters_without_videos:
        print("\nPosters without videos:")
        for name in posters_without_videos:
            print(f"- {name}")
    
    print("-" * 50)

@app.route('/')
def index():
    total_videos = sum(1 for data in video_data.values() if data['video'] and data['poster'])
    return render_template('index.html', total_videos=total_videos)

@app.route('/api/videos')
def get_videos():
    query = request.args.get('q', '').lower()
    filtered_videos = {}
    
    for base_name, data in video_data.items():
        if query in base_name.lower():
            if data['video'] and data['poster']:
                filtered_videos[base_name] = {
                    'title': base_name,
                    'video_path': f'/api/video/{base_name}',
                    'poster_path': f'/api/poster/{base_name}'
                }
    
    return jsonify(list(filtered_videos.values()))

@app.route('/api/video/<base_name>')
def get_video(base_name):
    if base_name in video_data and video_data[base_name]['video']:
        video_path = video_data[base_name]['video']
        mime_type = mimetypes.guess_type(video_path)[0] or 'video/mp4'
        return send_file(video_path, mimetype=mime_type)
    return '', 404

@app.route('/api/poster/<base_name>')
def get_poster(base_name):
    if base_name in video_data and video_data[base_name]['poster']:
        poster_path = video_data[base_name]['poster']
        mime_type = mimetypes.guess_type(poster_path)[0] or 'image/jpeg'
        return send_file(poster_path, mimetype=mime_type)
    return '', 404

if __name__ == '__main__':
    # Scan all configured directories on startup
    print("Starting video directory scan...")
    for directory in VIDEO_DIRECTORIES:
        scan_directory(directory)
    print_scan_summary()
    
    app.run(debug=True) 