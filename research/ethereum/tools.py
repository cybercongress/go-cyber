

def human_readable_time_interval(interval):
    h = int(interval / (60 * 60))
    m = int((interval % (60 * 60)) / 60)
    s = interval % 60.
    return "{}h: {:>02}m: {:>05.5f}s".format(h, m, s)
# End human_readable_time_interval
