import secrets

import bcrypt


def validate_password(password, password_hash):
    # bcrypt requires bytes, handle both string and bytes input
    if isinstance(password, str):
        password = password.encode('utf-8')
    if isinstance(password_hash, str):
        password_hash = password_hash.encode('utf-8')
    return bcrypt.checkpw(password, password_hash)


def generate_key():
    return secrets.token_urlsafe(21)
