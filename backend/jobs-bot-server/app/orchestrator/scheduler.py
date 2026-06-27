from apscheduler.schedulers.background import BackgroundScheduler

scheduler = BackgroundScheduler()

# Example stub for scheduling a run


def schedule_run(cron_expr, run_func):
    scheduler.add_job(run_func, 'cron', **cron_expr)


scheduler.start()
