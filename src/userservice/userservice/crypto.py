import bcrypt


def generate_hash(password):
    """
    Generate a hash for the given password.
    
    Args:
        password (str): The password to hash.
    
    Returns:
        str: The hashed password.
    """
    if isinstance(password, str):
        password = password.encode('utf-8')
    return bcrypt.hashpw(password, bcrypt.gensalt()).decode('utf-8')


def validate_password(password, password_hash):
    """
    Validate the given password against the hashed password.
    
    Args:
        password (str): The password to validate.
        password_hash (str): The hashed password.
    
    Returns:
        bool: True if the password matches the hash, False otherwise.
    """
    if isinstance(password, str):
        password = password.encode('utf-8')
    if isinstance(password_hash, str):
        password_hash = password_hash.encode('utf-8')
    return bcrypt.checkpw(password, password_hash)
