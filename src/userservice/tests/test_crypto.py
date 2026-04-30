import bcrypt
import pytest
from userservice.crypto import generate_hash, validate_password


class TestGenerateHash:
    def test_generates_hash_from_string(self):
        password = "mysecretpassword"
        result = generate_hash(password)
        assert isinstance(result, str)
        assert len(result) > 0
        # Verify it's a valid bcrypt hash by checking format
        assert result.startswith('$2b$')

    def test_generates_different_hashes_for_same_password(self):
        password = "mysecretpassword"
        hash1 = generate_hash(password)
        hash2 = generate_hash(password)
        # Same password should generate different hashes due to salt
        assert hash1 != hash2
        # But both should validate
        assert validate_password(password, hash1) is True
        assert validate_password(password, hash2) is True

    def test_generates_hash_from_bytes(self):
        password = b"mysecretpassword"
        result = generate_hash(password)
        assert isinstance(result, str)
        assert result.startswith('$2b$')


class TestValidatePassword:
    def test_valid_password_with_string_hash(self):
        password = "mysecretpassword"
        password_hash = generate_hash(password)
        assert validate_password(password, password_hash) is True

    def test_valid_password_with_bytes_password(self):
        password = b"mysecretpassword"
        password_hash = generate_hash("mysecretpassword")
        assert validate_password(password, password_hash) is True

    def test_invalid_password(self):
        password = "mysecretpassword"
        wrong_password = "wrongpassword"
        password_hash = generate_hash(password)
        assert validate_password(wrong_password, password_hash) is False

    def test_mixed_string_and_bytes(self):
        password = "mysecretpassword"
        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
        # hash is bytes, password is string
        assert validate_password(password, password_hash) is True

    def test_invalid_password_with_bytes(self):
        password = b"mysecretpassword"
        wrong_password = b"wrongpassword"
        password_hash = bcrypt.hashpw(password, bcrypt.gensalt())
        assert validate_password(wrong_password, password_hash) is False
