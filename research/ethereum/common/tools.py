import time
from _datetime import datetime
from statistics import mean
from typing import Dict

Rank = Dict[str, float]


def print_with_time(message: str):
    print("{} ".format(datetime.now().time()) + message)


def human_readable_interval(interval):
    h = int(interval / (60 * 60))
    m = int((interval % (60 * 60)) / 60)
    s = interval % 60.
    return "{}h: {:>02}m: {:>05.5f}s".format(h, m, s)


def diff(mx, mn):
    return abs(mx - mn)


def ranks_diff(a: Rank, b: Rank) -> (float, float):
    """
    Calculated mean and average ranks difference
    :param a: first rank
    :param b: second rank
    :return: (mean, average) ranks difference
    """
    diffs = []
    for node, a_node_rank in a.items():
        b_node_rank = b[node]
        rank_diff = diff(max(a_node_rank, b_node_rank), min(a_node_rank, b_node_rank))
        diffs.append(rank_diff)

    mean_rank_diff = mean(diffs)
    average_rank_diff = sum(diffs) / len(diffs)

    return mean_rank_diff, average_rank_diff


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
