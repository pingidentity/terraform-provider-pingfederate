import mmap, os
from glob import glob

# Verify that all resources contain an import.sh file and an entry in CHANGELOG.md
has_no_import_files = False
in_changelog = False
no_import_dirs = []
missing_changelog_entries = []

dirs = glob("./examples/resources/*/", recursive=True)
for dir in dirs:
  resource_name = str(dir).split("/")[-2].encode()
  with open(r'./CHANGELOG.md', 'rb', 0) as file:
    s = mmap.mmap(file.fileno(), 0, access=mmap.ACCESS_READ)
    if s.find(resource_name) != -1:
      in_changelog = True
      continue
    else:
      in_changelog = False
      if resource_name not in missing_changelog_entries:
        missing_changelog_entries.append(resource_name.decode())
      continue
  
  has_import_file = os.path.isfile(f"{dir}/import.sh")
  if not has_import_file:
    no_import_dirs.append(dir)
    has_no_import_files = True
  else:
    continue

if not in_changelog:
  print(f"Resource(s) {missing_changelog_entries} not found in CHANGELOG.md!")
  exit(1)
else:
  print("All resources found in CHANGELOG.md")

if has_no_import_files:
  print(f"No import.sh content found for resource(s) in {no_import_dirs}!") 
  exit(1)
else:
  print("All resources contain an import.sh for documentation")