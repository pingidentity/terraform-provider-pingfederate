import os
from glob import glob

dirs = glob("./examples/resources/*/", recursive=True)

has_no_import_files = False
no_import_dirs = []

for dir in dirs:
  has_import_file = os.path.isfile(f"{dir}/import.sh")
  if not has_import_file:
    no_import_dirs.append(dir)
    has_no_import_files = True
  else:
    continue

if has_no_import_files:
  print(f"No import.sh content found for resource(s) in {no_import_dirs}!") 
  exit(1)
else:
  print("All resources contain an import.sh for documentation")