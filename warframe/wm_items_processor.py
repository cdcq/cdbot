 
import json
import requests
import sqlite3

url = 'https://api.warframe.market/v1/items'
headers = {'accept': 'application', 'Language': 'zh-hans'}
s = requests.get(url=url, headers=headers).text

data = json.loads(s)

data = data['payload']['items']

for i in range(len(data)):
    data[i]['item_name'] = data[i]['item_name'].replace(' ', '')

conn = sqlite3.connect("./database.db")
curs = conn.cursor()

# nick_names = yaml.load(open('nick_names.yaml', 'r', encoding='utf-8').read(), Loader=yaml.FullLoader)

for i in data:
    item_name = i['item_name']
    # nick_name = ''
    # if 'Prime' in item_name and item_name[:item_name.find('Prime')] in nick_names:
    #     nick_name = nick_names[item_name[:item_name.find('Prime')]]

    curs.execute("INSERT INTO WM_ITEMS(ID, NAME, URL_NAME) VALUES(?, ?, ?)",
                 (i['id'], i['item_name'], i['url_name']))

conn.commit()
