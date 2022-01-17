import json
import requests
import getpass
from pprint import pprint
import os
import time
import datetime
import math

# curl "https://api.github.com/search/repositories?q=language:go&sort=stars&order=desc" > gorepo.json


def __header_to_dict(json_to_convert) -> dict:
    new_dict = {}
    for element in json_to_convert.keys():
        temp_dict = {}
        if type(element) is dict:
            new_dict[element] = __header_to_dict(
                json_to_convert[element])
        else:
            new_dict[element] = json_to_convert[element]
    return new_dict


i = 0
p = 0
cont = True

minimum_star_count = 2000

username = input("GitHub username: ")
pswd = getpass.getpass()

s = requests.Session()
s.auth = (username, pswd)
s.headers.update({'Accept': 'application/vnd.github.mercy-preview+json'})

visited_projects = []

keywords = []

while cont:

    p += 1
    payload = {'q': 'language:go', 'sort': 'stars',
               'order': 'desc', 'per_page': '100', 'page': str(p)}
    r = s.get('https://api.github.com/search/repositories', params=payload)
    data = r.json()
    # print(data)
    if 'items' not in data.keys():
        # check if API timeout
        if int(r.headers['X-RateLimit-Remaining']) == 0:
            print(
                f'\nwaiting for API limit reset, will try upto 5 times')
            attempts = 5
            while attempts > 0 and int(r.headers['X-RateLimit-Remaining']) == 0:
                attempt_now = math.floor(
                    datetime.datetime.utcnow().timestamp())
                reset = int(r.headers["X-RateLimit-Reset"])
                wait = reset - attempt_now + 3

                # print helpful for writing this
                print(
                    f'\n\tattempts: {attempts}\n\tlimit: {r.headers["X-RateLimit-Limit"]}\n\tremaining: {r.headers["X-RateLimit-Remaining"]}\n\treset: {reset}\n\twait: {wait} seconds')

                time.sleep(wait)
                # try again
                r = s.get(
                    'https://api.github.com/search/repositories', params=payload)
                data = r.json()
                attempts -= 1

            if 'items' in data.keys():
                print('\nsuccess waiting for API limit to reset\n')
            elif 'message' in data.keys():
                print(
                    f'\nWARNING: failure to wait for API limit to reset. last received request will be saved to \'error_output_result_limit.json\'\n\n\tmessage: {data["message"]}\n\tdocumentation_url: {data["documentation_url"]}\n')
                cont = False
                f = open(os.path.join(
                    os.getcwd(), 'error_output_result_limit.json'), 'w+')
                f.write(
                    f'{{"header":{json.dumps(__header_to_dict(r.headers), indent=4, sort_keys=True)},"data":{json.dumps(data, indent=4, sort_keys=True)}}}')
                f.close()
            else:
                print('\nERROR: failure due to unknown reasons. cannot continue. last received request will be saved to \'error_output_unknown_after.json\'\n')
                f = open(os.path.join(
                    os.getcwd(), 'error_output_unknown_after.json'), 'w+')
                f.write(
                    f'{{"header":{json.dumps(__header_to_dict(r.headers), indent=4, sort_keys=True)},"data":{json.dumps(data, indent=4, sort_keys=True)}}}')
                f.close()
                quit()
        else:
            print('\nERROR: failure due to unknown reasons. cannot continue. last received request will be saved to \'error_output_unknown.json\'\n')
            f = open(os.path.join(
                os.getcwd(), 'error_output_unknown.json'), 'w+')
            f.write(
                f'{{"header":{json.dumps(__header_to_dict(r.headers), indent=4, sort_keys=True)},"data":{json.dumps(data, indent=4, sort_keys=True)}}}')
            f.close()
            quit()
    else:
        for repo in data['items']:
            i += 1
            #   r = s.get('https://api.github.com/repos/'+str(repo['full_name']))
            #   repodata = r.json()

            if 'topics' not in repo.keys():
                print(f'\nERROR: topics was not found in repo (even if empty), something is wrong. cannot continue. last received request will be saved to \'error_output_repo.json\'\n')
                f = open(os.path.join(
                    os.getcwd(), 'error_output_repo.json'), 'w+')
                f.write(
                    f'{{"header":{json.dumps(__header_to_dict(r.headers), indent=4, sort_keys=True)},"repo_data":{json.dumps(repo, indent=4, sort_keys=True)}}}')
                f.close()
                quit()
            else:
                if len(repo['topics']) == 0:
                    print(str(repo['full_name']), ','                      # , str(repo['description']),','
                          , str(repo['watchers_count']))
                else:
                    print(str(repo['full_name']), ','                      # , str(repo['description']),','
                          , str(repo['watchers_count']), ',', ','.join(repo['topics']))
                    # keywords = keywords+(repodata['topics'])

                visited_projects.append(str(repo['full_name']))

            if int(str(repo['watchers_count'])) < minimum_star_count:
                cont = False

print(str(i), " projects found")

f = open(os.path.join(os.getcwd(), 'projects.txt'), 'w+')
for p in visited_projects:
    f.write(f'{p}\n')
f.close()
print('written to \'projects.txt\'')
# print(','.join(set(keywords)))
