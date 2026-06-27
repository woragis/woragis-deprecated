import datetime


def now_utc():
    return datetime.datetime.utcnow()


def format_time(dt):
    return dt.strftime('%Y-%m-%d %H:%M:%S')
