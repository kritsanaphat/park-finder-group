
def check_exists_licence_plate(licence_plate,json_data):
    for i in json_data:
        print(i["lpr"])
        if i["lpr"] == licence_plate:
            return True
        
    return False