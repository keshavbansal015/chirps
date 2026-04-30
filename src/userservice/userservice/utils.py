import re


def is_valid_email(email):
    # Basic email validation: no spaces allowed anywhere
    return re.match(r'[^\s@]+@[^\s@]+\.[^\s@]+', email) is not None
