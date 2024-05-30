import mmap, re, sys, os
from glob import glob

files_to_format = sys.argv[1:]
print(f"files needed for markdown formatting: {files_to_format}")

for file in files_to_format:
  file_to_format = f"./docs/resources/{file}"
  if os.path.exists(file_to_format):
    print(f"formatting {file_to_format}...")
    with open(rf'{file_to_format}', 'r+') as file_to_edit:
      file_content = file_to_edit.read()
      file_content = re.sub('(## Schema\n\n```)','## Schema\n', file_content)
      file_content = re.sub('(```\n\n## Import)','\n## Import', file_content)

      # Setting the position to the beginning of the file
      file_to_edit.seek(0) 
        
      # Writing replaced data in the file 
      file_to_edit.write(file_content) 

      # Truncating the file size 
      file_to_edit.truncate() 
      print(f"{file_to_format} formatted successfully!")