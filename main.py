import sys
import os

print(f"Python {sys.version} running in {sys.platform}/wazero.")

for root, dirs, files in os.walk("/", topdown=False):
   for name in dirs:
      print(f"python fs file: {os.path.join(root, name)}/")
   for name in files:
      print(f"python fs file: {os.path.join(root, name)}")

try:
   with open('/output/from-python.txt', 'w') as f:
      f.write('from python')
except Exception as e:
   print(f'ERROR: failed to write /output/from-python.txt: {e}')
