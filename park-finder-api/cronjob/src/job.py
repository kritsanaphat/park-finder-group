from datetime import datetime, timedelta
from http import HTTPStatus
from apscheduler.triggers.date import DateTrigger


from flask import Blueprint, request


job_api = Blueprint("job_api", __name__, url_prefix="/job")

@job_api.route("/", methods=["Get"])
def hello():
    return hello

# สำหรับการยกเลิกการจอง เมื่อไม่จ่ายเงินตามกำหนด
@job_api.route("/cancel_reserve", methods=["POST"]) #CR
def add_cancel_reserve_job(order_id="",time_trigger = ""):
    from . import scheduler, status_logs, scheduler_logs,tz

    if order_id == "":
        order_id = request.args.get("order_id")

    try:
        start_time = get_run_time(20)
        job_id = str("CR_"+order_id)

        if time_trigger == "":
            trigger = DateTrigger(run_date=start_time)
        else:
            trigger = DateTrigger(run_date=create_trigger(time_trigger))
        
        scheduler_logs.save_job(job_id,str(start_time))
        scheduler.add_job(
            func=status_logs.CR_cancel_reserve,
            trigger=trigger,
            id=job_id,
            kwargs={'order_id': order_id}
        )
        print("Job scheduled cancel reserve successfully, Job id :"+job_id)
        return {"status": 200, "message": "Job scheduled cancel reserve successfully, Job id :"+job_id}

    except Exception as e:
        print("Exception at add_cancel_reserve_job function:"+str(e))
        return {"status": 400, "message": "Exception at add_cancel_reserve_job function:"+str(e)}

# สำหรับการยกเลิกการจอง เมื่อไม่จ่ายเงินตามกำหนด
@job_api.route("/cancel_extend_reserve", methods=["POST"]) #CR
def add_cancel_extend_reserve_job(order_id="",time_trigger = ""):
    from . import scheduler, status_logs, scheduler_logs,tz

    if order_id == "":
        order_id = request.args.get("order_id")

    try:
        start_time = get_run_time(20)
        job_id = str("CER_"+order_id)

        if time_trigger == "":
            trigger = DateTrigger(run_date=start_time)
            scheduler_logs.save_job(job_id,str(start_time))
        else:
            trigger = DateTrigger(run_date=create_trigger(time_trigger))

        scheduler.add_job(
            func=status_logs.CER_cancel_extend_reserve,
            trigger=trigger,
            id=job_id,
            kwargs={'order_id': order_id}
        )
        print("Job scheduled cancel extend reserve successfully, Job id :"+job_id)
        return {"status": 200, "message": "Job scheduled cancel extend reserve successfully, Job id :"+job_id}

    except Exception as e:
        print("Exception at add_cancel_reserve_job function:"+str(e))
        return {"status": 400, "message": "Exception at add_cancel_reserve_job function:"+str(e)}

# สำหรับการยกเลิกการจองล่วงหน้า เมื่อเจ้าของที่ไม่กดยืนยันในเวลาที่กำหนด
@job_api.route("/cancel_reserve_in_advance", methods=["POST"]) #CR
def add_cancel_reserve_job_in_advance(order_id="",time_trigger = ""):
    from . import scheduler, status_logs, scheduler_logs,tz

    if order_id == "":
        order_id = request.args.get("order_id")

    try:
        order_id=order_id.split(",")
        start_time = get_run_time(1440)
        job_id = str("CRIA_"+order_id[0]+","+order_id[1])

        if time_trigger == "":
            trigger = DateTrigger(run_date=start_time)
            scheduler_logs.save_job(job_id,str(start_time))
        else:
            trigger = DateTrigger(run_date=create_trigger(time_trigger))

        # Add the job to the scheduler
        scheduler.add_job(
            func=status_logs.CRIA_cancel_reserve_in_advance,
            trigger=trigger,
            id=job_id,
            kwargs={'order_id': order_id}
        )
        return {"status": 200, "message": "Job scheduled cancel reserve in advancesuccessfully, Job id :"+job_id}

    except Exception as e:
        return {"status": 400, "message": "Exception at add_cancel_reserve_job function:"+str(e)}

# สำหรับแจ้งเตือนเมื่อหมดเวลาที่จอง
@job_api.route("/timeout_reserve", methods=["POST"]) #TOR
def add_timeout_reserve_job(order_id="",time_trigger="",BTOR=True,TOR=True,ATOR=True):
    from . import scheduler, status_logs, scheduler_logs,tz
    #from fetch
    if order_id == "":
        order_id = request.args.get("order_id")
        date = request.args.get("date")
        hour = request.args.get("hour")
        min = request.args.get("min")
        date = date.split("-")
        start_time = get_run_time_by_date(
                    int(date[0]), int(date[1]), int(date[2]), int(hour), int(min)
                )
            #before 15 minutes 
        if BTOR:
            try:
                job_id = str("BTOR_" + order_id)
                # Remove 15 minutes to the start_time
                start_time -= timedelta(minutes=15)
                trigger = DateTrigger(run_date=start_time)
                scheduler_logs.save_job(job_id, str(start_time))
                scheduler.add_job(
                    func=status_logs.BTOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                return {"status": 400, "message": "Exception at add_before_timeout_reserve function:"+str(e)}
        
        #on time 
        if TOR:
            try:
                job_id = str("TOR_"+order_id)
                start_time += timedelta(minutes=15)
                trigger = DateTrigger(run_date=start_time)
                scheduler_logs.save_job(job_id,str(start_time))
                scheduler.add_job(
                    func=status_logs.TOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                return {"status": 400, "message": "Exception at add_timeout_reserve function:"+str(e)}
        
        ## after 15 minutes
        if ATOR:
            try:
                job_id = str("ATOR_" + order_id)
                # Add 15 minutes to the start_time
                start_time += timedelta(minutes=15)
                trigger = DateTrigger(run_date=start_time)
                scheduler_logs.save_job(job_id, str(start_time))
                scheduler.add_job(
                    func=status_logs.ATOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                return {"status": 400, "message": "Exception at add_after_timeout_reserve function:"+str(e)}
    
    else:
        print("Case Fetch")
            #before 15 minutes 
        if BTOR:          
            try:
                job_id = str("BTOR_" + order_id)
                trigger = DateTrigger(run_date=create_trigger(time_trigger))
                scheduler_logs.save_job(job_id, str(time_trigger))
                scheduler.add_job(
                    func=status_logs.BTOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                print("Exception at add_before_timeout_reserve function:"+str(e))
                return {"status": 400, "message": "Exception at add_before_timeout_reserve function:"+str(e)}
        
        #on time 
        if TOR:
            try:
                job_id = str("TOR_"+order_id)
                trigger = DateTrigger(run_date=create_trigger(time_trigger))
                scheduler_logs.save_job(job_id, str(time_trigger))
                scheduler.add_job(
                    func=status_logs.TOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                print("Exception at add_timeout_reserve function:"+str(e))
                return {"status": 400, "message": "Exception at add_timeout_reserve function:"+str(e)}
        
        ## after 15 minutes
        if ATOR:
            try:
                job_id = str("ATOR_" + order_id)
                trigger = DateTrigger(run_date=create_trigger(time_trigger))
                scheduler_logs.save_job(job_id, str(time_trigger))
                scheduler.add_job(
                    func=status_logs.ATOR_timeout_reserve_job,
                    trigger=trigger,
                    id=job_id,
                    kwargs={'order_id': order_id}
                )
            except Exception as e:
                print("Exception at add_after_timeout_reserve function:"+str(e))
                return {"status": 400, "message": "Exception at add_after_timeout_reserve function:"+str(e)}

    return {"status": 200, "message": "Job scheduled timeout reserve successfully, Job id :"+job_id}

# สำหรับการอัพเดทเวลาปิด-เปิดที่จอดรถที่ตั้งไว้
@job_api.route("/update_open_area_status", methods=["POST"]) #OPA
def add_update_area_open_status_job(parking_id="",time_trigger = ""):
    from . import scheduler, parking_area_logs, scheduler_logs,tz

    if parking_id == "":
        parking_id = request.args.get("parking_id")

    try:
        range_time = request.args.get("range_time")
        start_time = get_run_time(int(range_time))
        job_id = str("OPA"+parking_id)

        if time_trigger == "":
            trigger = DateTrigger(run_date=start_time)
            scheduler_logs.save_job(job_id,str(start_time))
        else:
            trigger = DateTrigger(run_date=create_trigger(time_trigger))

        scheduler.add_job(
            func=parking_area_logs.update_open_area_status,
            trigger=trigger,
            id=job_id,
            kwargs={'parking_id': parking_id}
        )
        print("Job scheduled update open area status successfully, Job id :"+job_id)
        return {"status": 200, "message": "Job scheduled add update area open status job successfully, Job id :"+job_id}

    except Exception as e:
        print("Exception at add_cancel_reserve_job function:"+str(e))
        return {"status": 400, "message": "Exception at add_update_area_open_status_job function:"+str(e)}

@job_api.route("/remove_job", methods=["POST"])
def remove_job():
    from . import scheduler, status_logs, scheduler_logs,tz

    job_id = request.args.get("job_id")
    try:
        if get_job(job_id) == None:
            return {"status": 200, "message": "Not Found, Job id :"+job_id}
        scheduler.remove_job(
            job_id=job_id
        )

        print("Remove job id ", id)
        scheduler_logs.remove_job(job_id)
        return {"status": 200, "message": "Remove scheduled successfully, Job id :"+job_id}

    except Exception as e:
        return {"status": 400, "message": (f"Exception at remove_job function: {e}")}
    

def remove_job_fetch(job_id):
    from . import scheduler, status_logs, scheduler_logs,tz

    try:
        if get_job(job_id) == None:
            return {"status": 200, "message": "Not Found, Job id :"+job_id}
        scheduler.remove_job(
            job_id=job_id
        )

        print("Remove job id ", id)
        scheduler_logs.remove_job(job_id)
        return {"status": 200, "message": "Remove scheduled successfully, Job id :"+job_id}

    except Exception as e:
        return {"status": 400, "message": (f"Exception at remove_job function: {e}")}

def reschedule_chat_timeout_job(job_id: str, timeout: int):
    from . import scheduler, status_logs, scheduler_logs,tz

    try:
        scheduler.reschedule_job(
            job_id=job_id,
            trigger="date",
            run_date=get_run_time(timeout)
        )
    except Exception as e:
        print(f"Exception at reschedule_job function: {e}")

def get_job(job_id: str):
    from . import scheduler

    return scheduler.get_job(
        job_id=job_id
    )

def get_run_time(minutes: int) -> datetime:
    from . import tz
    datetime_now = datetime.now(tz=tz)
    return datetime_now + timedelta(minutes=minutes)

def get_run_time_by_date(year, month, day, hour, min):
    from . import tz
    return tz.localize(datetime(year, month, day, hour, min))

def fetch_scheduler_logs(job_list):
    for job_dict in job_list:
        for key, value in job_dict.items():
            remove_job_fetch(key)
            if key[:2]=="CR" and key[:4]!="CRIA":
                add_cancel_reserve_job(key[3:],value)
            elif key[:3]=="CER":
                add_cancel_extend_reserve_job(key[4:],value)
            elif key[:4]=="BTOR":
                add_timeout_reserve_job(key[5:],value,True,False,False)
            elif key[:3]=="TOR":
                add_timeout_reserve_job(key[4:],value,False,True,False)
            elif key[:4]=="ATOR":
                add_timeout_reserve_job(key[5:],value,False,False,True)
            elif key[:4]=="CRIA":
                add_cancel_reserve_job_in_advance(key[5:],value)



def create_trigger(time_str):
    from datetime import datetime, timedelta, timezone
    time_trigger_datetime = datetime.fromisoformat(time_str.replace(' ', 'T'))
    tz_offset = int(time_trigger_datetime.utcoffset().total_seconds() / 60)
    tz_offset = timedelta(minutes=tz_offset)
    tzinfo = timezone(tz_offset)
    time_trigger_datetime = time_trigger_datetime.replace(tzinfo=tzinfo)

    return time_trigger_datetime