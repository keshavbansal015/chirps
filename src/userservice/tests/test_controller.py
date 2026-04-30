import pytest
from unittest.mock import Mock
from grpc import StatusCode

from userservice.controller import Controller
from userservice import chirp_pb2


class MockContext:
    """Mock gRPC context for testing"""
    def __init__(self, user_id='123'):
        self.abort_called = False
        self.abort_code = None
        self.abort_details = None
        self._metadata = {'user-id': user_id}

    def invocation_metadata(self):
        return [('user-id', self._metadata['user-id'])]

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


class TestCreateUser:
    def test_create_user_success(self, controller, mock_db_client, mock_context):
        mock_db_client.create_user.return_value = 1

        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='johndoe',
            email='john@example.com',
            password='password123'
        )

        result = controller.CreateUser(request, mock_context)

        assert result.id == 1
        mock_db_client.create_user.assert_called_once()

    def test_create_user_missing_username(self, controller, mock_context):
        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='',
            email='john@example.com',
            password='password123'
        )

        with pytest.raises(Exception):
            controller.CreateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_user_missing_email(self, controller, mock_context):
        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='johndoe',
            email='',
            password='password123'
        )

        with pytest.raises(Exception):
            controller.CreateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_user_missing_password(self, controller, mock_context):
        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='johndoe',
            email='john@example.com',
            password=''
        )

        with pytest.raises(Exception):
            controller.CreateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_user_invalid_email(self, controller, mock_context):
        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='johndoe',
            email='invalid-email',
            password='password123'
        )

        with pytest.raises(Exception):
            controller.CreateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_create_user_duplicate(self, controller, mock_db_client, mock_context):
        mock_db_client.create_user.side_effect = Exception('Duplicate')

        request = chirp_pb2.CreateUserRequest(
            name='John Doe',
            username='johndoe',
            email='john@example.com',
            password='password123'
        )

        with pytest.raises(Exception):
            controller.CreateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT


class TestGetUser:
    def test_get_user_success(self, controller, mock_db_client, mock_context):
        mock_db_client.get_user.return_value = chirp_pb2.User(
            id=1, name='John', username='johndoe', email='john@example.com'
        )

        request = chirp_pb2.UserRequest(user_id=1)
        result = controller.GetUser(request, mock_context)

        assert result.id == 1
        assert result.username == 'johndoe'

    def test_get_user_not_found(self, controller, mock_db_client, mock_context):
        mock_db_client.get_user.return_value = None

        request = chirp_pb2.UserRequest(user_id=999)

        with pytest.raises(Exception):
            controller.GetUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.NOT_FOUND

    def test_get_user_db_error(self, controller, mock_db_client, mock_context):
        mock_db_client.get_user.side_effect = Exception('DB error')

        request = chirp_pb2.UserRequest(user_id=1)

        with pytest.raises(Exception):
            controller.GetUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INTERNAL


class TestUpdateUser:
    def test_update_user_success(self, controller, mock_db_client, mock_context):
        mock_db_client.update_user.return_value = None

        request = chirp_pb2.UpdateUserRequest(
            name='John Updated',
            username='johnupdated',
            email='john.updated@example.com',
            bio='Updated bio'
        )

        result = controller.UpdateUser(request, mock_context)

        assert isinstance(result, chirp_pb2.Empty)

    def test_update_user_invalid_email(self, controller, mock_context):
        request = chirp_pb2.UpdateUserRequest(
            name='John',
            username='johndoe',
            email='invalid-email',
            bio='Bio'
        )

        with pytest.raises(Exception):
            controller.UpdateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT


class TestUpdatePassword:
    def test_update_password_success(self, controller, mock_db_client, mock_context):
        import bcrypt
        password_hash = bcrypt.hashpw(b'oldpassword', bcrypt.gensalt()).decode('utf-8')
        mock_db_client.get_user_with_id.return_value = {'id': 123, 'password': password_hash}
        mock_db_client.update_password.return_value = None

        request = chirp_pb2.UpdateUserRequest(
            password='newpassword',
            old_password='oldpassword'
        )

        result = controller.UpdateUser(request, mock_context)

        assert isinstance(result, chirp_pb2.Empty)

    def test_update_password_missing_old_password(self, controller, mock_context):
        request = chirp_pb2.UpdateUserRequest(
            password='newpassword',
            old_password=''
        )

        with pytest.raises(Exception):
            controller.UpdateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT

    def test_update_password_wrong_old_password(self, controller, mock_db_client, mock_context):
        import bcrypt
        password_hash = bcrypt.hashpw(b'correctpassword', bcrypt.gensalt()).decode('utf-8')
        mock_db_client.get_user_with_id.return_value = {'id': 123, 'password': password_hash}

        request = chirp_pb2.UpdateUserRequest(
            password='newpassword',
            old_password='wrongpassword'
        )

        with pytest.raises(Exception):
            controller.UpdateUser(request, mock_context)

        assert mock_context.abort_called
        assert mock_context.abort_code == StatusCode.INVALID_ARGUMENT


class TestGetFollowing:
    def test_get_following_success(self, controller, mock_db_client, mock_context):
        mock_db_client.get_following.return_value = [
            chirp_pb2.User(id=2, username='user2'),
            chirp_pb2.User(id=3, username='user3')
        ]

        request = chirp_pb2.GetUsersRequest(user_id=1, page=0, limit=10)
        result = controller.GetFollowing(request, mock_context)

        assert len(result.users) == 2
        assert result.users[0].username == 'user2'


class TestGetFollowers:
    def test_get_followers_success(self, controller, mock_db_client, mock_context):
        mock_db_client.get_followers.return_value = [
            chirp_pb2.User(id=4, username='follower1')
        ]

        request = chirp_pb2.GetUsersRequest(user_id=1, page=0, limit=10)
        result = controller.GetFollowers(request, mock_context)

        assert len(result.users) == 1
        assert result.users[0].username == 'follower1'


class TestFollowUser:
    def test_follow_user_success(self, controller, mock_db_client, mock_context):
        mock_db_client.follow_user.return_value = None

        request = chirp_pb2.UserRequest(user_id=2)
        result = controller.FollowUser(request, mock_context)

        assert isinstance(result, chirp_pb2.Empty)
        mock_db_client.follow_user.assert_called_once_with(2, '123')


class TestUnfollowUser:
    def test_unfollow_user_success(self, controller, mock_db_client, mock_context):
        mock_db_client.unfollow_user.return_value = None

        request = chirp_pb2.UserRequest(user_id=2)
        result = controller.UnfollowUser(request, mock_context)

        assert isinstance(result, chirp_pb2.Empty)
        mock_db_client.unfollow_user.assert_called_once_with(2, '123')
