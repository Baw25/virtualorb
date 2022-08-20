from random import randrange
import base64
import os 
import requests
import sys
import time 
import names
import pprint

# Read the test iris images into encoded files
VIRTUAL_ORB_URL = "http://localhost:5000/v1/virtualorb"
ITERATIONS = "iterations"
INTERVAL = "interval"
SIGNUP_REQUEST = "signup"
REPORT_REQUEST = "report"

def encoded_input_images():
    encoded_images = []
    try: 
        for filename in os.listdir(os.path.join(os.getcwd(), "images")):
            with open(os.path.join(os.getcwd(), "images", filename), 'rb') as image_str: 
                current_encoded_string = base64.b64encode(image_str.read())
                encoded_images.append(current_encoded_string.decode("utf-8"))
    except IOError:
        print("Problem loading encoded images!")
    finally:
        return encoded_images 

def send_multiple_signup_requests():
    print("Sending signup request...")
    encoded_images = encoded_input_images()
    random_user_name = names.get_full_name()
    response = requests.post(VIRTUAL_ORB_URL + "/signup", json={"images": encoded_images, "name": random_user_name })
    json_response = response.json()
    print("Signup generated:")
    pprint.pprint(json_response)

def send_single_signup_request():
    print("Sending signup request...")
    encoded_images = encoded_input_images()
    random_index = randrange(4)
    image = encoded_images[random_index]
    random_user_name = names.get_full_name()

    response = requests.post(VIRTUAL_ORB_URL + "/signup", json={"images": [image], "name": random_user_name })
    json_response = response.json()
    print("Signup generated:")
    pprint.pprint(json_response)

def send_single_report_request():
    print("Sending status report request...")
    response = requests.post(VIRTUAL_ORB_URL + "/report", json={ "send_report": True })
    json_response = response.json()
    print("Report sent:")
    pprint.pprint(json_response)
    
def main():
    print("Running simulator")
    # If defined (iterations and interval), send either or both requests over time intervals for input iterations
    if len(sys.argv) > 1:
        iterations, interval = 1, 0
        if ITERATIONS in sys.argv[1]:
            iterations = int(sys.argv[1].split("=")[1])
        else:
            print("Please 1st argument properly ex) iterations=3")
            sys.exit() 

        if INTERVAL in sys.argv[2]:
            interval = int(sys.argv[2].split("=")[1])
        else:
            print("Please 2nd argument properly ex) interval=3")
            sys.exit() 

        current = 0
        if len(sys.argv) > 3 and REPORT_REQUEST in sys.argv[3]:
            while current < iterations:
                send_single_report_request()
                current += 1 
                if current < iterations:
                    time.sleep(interval)  
        else:
            while current < iterations:
                send_single_signup_request()
                current += 1 
                if current < iterations:
                    time.sleep(interval)                                

    # Just send each of the requests (report and signup) once
    else:
        send_multiple_signup_requests()
        time.sleep(2) 
        send_single_report_request()

if __name__ == "__main__":
    main()