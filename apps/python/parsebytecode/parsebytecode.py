from evmole import contract_info
import json 
import re
import sys

code= sys.argv[1] 

info = contract_info(code, selectors=True, arguments=True, state_mutability=True, storage=True)

functions = []
storage = []



for func in info.functions:
    functions.append(
        {
            'selector': func.selector, 
            'arguments': func.arguments, 
            'state_mutability': func.state_mutability,
            'bytecode_offset': func.bytecode_offset
        }
    )

for element in info.storage:
    
    storage_record_str = str(element)

    # Regular expression pattern to extract the values
    pattern = r'StorageRecord\(slot=(.*?),offset=(.*?),type=(.*?),reads=\[(.*?)\],writes=\[(.*?)\]\)'

    # Match the pattern
    match = re.match(pattern, storage_record_str)

    if match:
        slot = match.group(1)
        offset = int(match.group(2))
        type = match.group(3)
        reads = match.group(4).replace('"', '').split(', ')
        writes = match.group(5).replace('"', '').split(', ')

        storage.append(
        {
            'slot': slot,
            'offset': offset,
            'type': type,
            'reads': reads,
            'writes': writes
        }
        )


info_dict = {
    'functions': functions,
    'storage': storage
}
jsonstr1 = json.dumps(info_dict, default=str)
print(jsonstr1)
