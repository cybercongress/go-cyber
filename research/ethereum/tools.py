import time
from _datetime import datetime


def print_with_time(message: str):
    print("{} ".format(datetime.now().time()) + message)


def human_readable_interval(interval):
    h = int(interval / (60 * 60))
    m = int((interval % (60 * 60)) / 60)
    s = interval % 60.
    return "{}h: {:>02}m: {:>05.5f}s".format(h, m, s)


def record_execution_time(before_message: str = None, after_message: str = None):
    """
    Decorator. Records and prints function execution time.
    :param before_message: string, printed before function execution
    :param after_message: string, containing single '{}', that will be replaced  by execution time
    :return: function execution result
    """

    def decorator(func):
        def wrapped(*args, **kwargs):

            if before_message:
                print(before_message)

            start_time = time.time()
            result = func(*args, **kwargs)
            end_time = time.time()
            readable_interval = human_readable_interval(end_time - start_time)

            if after_message:
                print(after_message.format(readable_interval))
            else:
                print("{} was executed in {}".format(func.__name__, readable_interval))

            return result

        return wrapped

    return decorator
