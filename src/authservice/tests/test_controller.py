import bcrypt
import pytest
from unittest.mock import Mock, MagicMock
from grpc import StatusCode

from authservice.controller import Controller
from authservice import chirp_pb2


class MockContext:
    """Mock gRPC context for testing"""
    def __init__(self):
        self.abort_called = False
        self.abort_code = None
        self.abort_details = None

    def abort(self, code, details=''):
        self.abort_called = True
        self.abort_code = code
        self.abort_details = details
        raise Exception(f'gRPC abort: {code} - {details}')


@pytest.fixture
def mock_db_client():
    return Mock()


@pytest.fixture
def controller(mock_db_client):
    return Controller(mock_db_client)


@pytest.fixture
def mock_context():
    return MockContext()


class TestCreateSession:
    def test_create_session_success(self, controller, mock_db_client, mock_context):
        # Setup - use a properly hashed password
        password_hash = bcrypt.hashpw(b'password', bcrypt.gensalt()).decode('utf-8')
        mock_db_client.get_user.return_value = {'id': 1, 'password': password_hash}
        mock_db_client.create_session.return_value = chirp_pb2.Session(
            id='session-123', user_id=1, created='2024-01-01'
        )

        request = chirp_pb2.Credentials(email='test@example.com', password='password')

        # Execute
        result = controller.CreateSession(request, mock_context)

        # Assert
        assert result.id == 'session-123'
        assert result.user_id == 1
        mock_db_client.get_user.assert_called_once_with('test@example.com')

    def test_create_session_missing_email(self, controller, mock_context):
        request = chirp_pb2.Credentials(email='', password='password')

        with pytest.raises(Exception):
            controller.CreateSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_session_missing_password(self, controller, mock_context):
        request = chirp_pb2.Credentials(email='test@example.com', password='')

        with pytest.raises(Exception):
            controller.CreateSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_session_user_not_found(self, controller, mock_db_client, mock_context):
        mock_db_client.get_user.return_value = None

        request = chirp_pb2.Credentials(email='test@example.com', password='password')

        with pytest.raises(Exception):
            controller.CreateSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.UNAUTHENTICATED

    def test_create_session_wrong_password(self, controller, mock_db_client, mock_context):
        import bcrypt
        hashed = bcrypt.hashpw(b'correct', bcrypt.gensalt()).decode('utf-8')
        mock_db_client.get_user.return_value = {'id': 1, 'password': hashed}

        request = chirp_pb2.Credentials(email='test@example.com', password='wrong')

        with pytest.raises(Exception):
            controller.CreateSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.UNAUTHENTICATED

    def test_create_session_db_error(self, controller, mock_db_client, mock_context):
        mock_db_client.get_user.side_effect = Exception('DB error')

        request = chirp_pb2.Credentials(email='test@example.com', password='password')

        with pytest.raises(Exception):
            controller.CreateSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INTERNAL


class TestGetSession:
    def test_get_session_success(self, controller, mock_db_client, mock_context):
        mock_db_client.get_session.return_value = chirp_pb2.Session(
            id='session-123', user_id=1, created='2024-01-01'
        )

        request = chirp_pb2.SessionRequest(session_id='session-123')
        result = controller.GetSession(request, mock_context)

        assert result.id == 'session-123'
        assert result.user_id == 1

    def test_get_session_not_found(self, controller, mock_db_client, mock_context):
        mock_db_client.get_session.return_value = None

        request = chirp_pb2.SessionRequest(session_id='invalid-session')

        with pytest.raises(Exception):
            controller.GetSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.UNAUTHENTICATED

    def test_get_session_db_error(self, controller, mock_db_client, mock_context):
        mock_db_client.get_session.side_effect = Exception('DB error')

        request = chirp_pb2.SessionRequest(session_id='session-123')

        with pytest.raises(Exception):
            controller.GetSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INTERNAL


class TestDeleteSession:
    def test_delete_session_success(self, controller, mock_db_client, mock_context):
        mock_db_client.delete_session.return_value = None

        request = chirp_pb2.SessionRequest(session_id='session-123')
        result = controller.DeleteSession(request, mock_context)

        assert isinstance(result, chirp_pb2.Empty)
        mock_db_client.delete_session.assert_called_once_with('session-123')

    def test_delete_session_db_error(self, controller, mock_db_client, mock_context):
        mock_db_client.delete_session.side_effect = Exception('DB error')

        request = chirp_pb2.SessionRequest(session_id='session-123')

        with pytest.raises(Exception):
            controller.DeleteSession(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INTERNAL
