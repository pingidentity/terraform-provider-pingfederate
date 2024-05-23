import mmap, re, sys
from glob import glob

files_to_format = sys.argv[1:]
print(f"files needed for markdown formatting: {files_to_format}")

resource_docs_location = glob("./docs/resources/*", recursive=True)
for file in files_to_format:
  for mdFile in resource_docs_location:
    if f"./docs/resources/{file}" == mdFile:
      print(f"formatting {mdFile}...")
      with open(rf'{mdFile}', 'r+') as file_to_edit:
        file_content = file_to_edit.read()
        file_content = re.sub('(## Schema\n\n```)','## Schema\n', file_content)
        file_content = re.sub('(```\n\n## Import)','\n## Import', file_content)

        # Setting the position to the beginning of the file
        file_to_edit.seek(0) 
          
        # Writing replaced data in the file 
        file_to_edit.write(file_content) 
  
        # Truncating the file size 
        file_to_edit.truncate() 
        print(f"{mdFile} formatted successfully!")