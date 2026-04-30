import bcrypt
import pytest
from authservice.crypto import validate_password, generate_key


class TestValidatePassword:
    def test_valid_password_with_bytes(self):
        password = b"mysecretpassword"
        password_hash = bcrypt.hashpw(password, bcrypt.gensalt())
        assert validate_password(password, password_hash) is True

    def test_valid_password_with_strings(self):
        password = "mysecretpassword"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
        assert validate_password(password, password_hash) is True

    def test_invalid_password_with_bytes(self):
        password = b"mysecretpassword"
        wrong_password = b"wrongpassword"
        password_hash = bcrypt.hashpw(password, bcrypt.gensalt())
        assert validate_password(wrong_password, password_hash) is False

    def test_invalid_password_with_strings(self):
        password = "mysecretpassword"
        wrong_password = "wrongpassword"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
        assert validate_password(wrong_password, password_hash) is False

    def test_mixed_string_and_bytes(self):
        password = "mysecretpassword"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
        assert validate_password(password, password_hash) is True


class TestGenerateKey:
    def test_generates_non_empty_string(self):
        key = generate_key()
        assert isinstance(key, str)
        assert len(key) > 0

    def test_generates_unique_keys(self):
        key1 = generate_key()
        key2 = generate_key()
        assert key1 != key2

    def test_generates_urlsafe_key(self):
        key = generate_key()
        # URL-safe characters are alphanumeric, hyphen, and underscore
        assert all(c.isalnum() or c in '-_' for c in key)
