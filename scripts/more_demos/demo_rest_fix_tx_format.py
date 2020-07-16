import json

with open('tx.json', 'r') as json_file:
    data = json.load(json_file)

new_data = dict()
new_data['tx'] = data['value']
new_data['mode'] = 'block'

with open('tx.json', 'w') as json_file:
    json.dump(new_data, json_file)
