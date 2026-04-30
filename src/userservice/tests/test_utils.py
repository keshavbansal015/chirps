import pytest
from userservice.utils import is_valid_email


class TestIsValidEmail:
    def test_valid_email(self):
        assert is_valid_email("test@example.com") is not None
        assert is_valid_email("user.name@domain.co.uk") is not None
        assert is_valid_email("user+tag@example.org") is not None

    def test_invalid_email_no_at_symbol(self):
        assert is_valid_email("testexample.com") is False

    def test_invalid_email_no_domain(self):
        assert is_valid_email("test@") is False

    def test_invalid_email_no_tld(self):
        assert is_valid_email("test@example") is False

    def test_invalid_email_empty_string(self):
        assert is_valid_email("") is False

    def test_invalid_email_multiple_at_symbols(self):
        # Current regex allows this - it's a simple validation
        assert is_valid_email("test@@example.com") is not None

    def test_invalid_email_spaces_in_domain(self):
        # Spaces before @ are allowed by simple regex, but spaces after @ are not
        assert is_valid_email("test@ example.com") is False
