import mmap, os
from glob import glob

# Verify that all resources contain an import.sh file and an entry in CHANGELOG.md
has_no_import_files = False
missing_in_changelog = False
no_import_dirs = set()
missing_changelog_entries = set()

noImportResources = [
  "pingfederate_connection_metadata_export",
  "pingfederate_keypairs_signing_key",
  "pingfederate_keypairs_ssl_server_csr",
  "pingfederate_keypairs_ssl_server_csr_export",
  "pingfederate_keypairs_ssl_server_key",
  "pingfederate_license"
]

dirs = glob("./examples/resources/*/", recursive=True)
for dir in dirs:
  resource_name = str(dir).split("/")[-2].encode()
  with open(r'./CHANGELOG.md', 'rb', 0) as file:
    s = mmap.mmap(file.fileno(), 0, access=mmap.ACCESS_READ)
    if s.find(resource_name) == -1:
      missing_changelog_entries.add(resource_name.decode())
  
  has_import_file = os.path.isfile(f"{dir}/import.sh")
  if not has_import_file and resource_name.decode() not in noImportResources:
    no_import_dirs.add(dir)

if len(missing_changelog_entries) > 0:
  missing_in_changelog = True
  print(f"Resource(s) {missing_changelog_entries} not found in CHANGELOG.md!")
else:
  print("All resources found in CHANGELOG.md")

if len(no_import_dirs) > 0:
  has_no_import_files = True
  print(f"No import.sh content found for resource(s) in {no_import_dirs}!") 
else:
  print("All resources contain an import.sh for documentation")

if has_no_import_files or missing_in_changelog:
  exit(1)