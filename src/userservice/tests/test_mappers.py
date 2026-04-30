from userservice.mappers import map_user
from userservice import chirp_pb2


class TestMapUser:
    def test_map_user_with_valid_data(self):
        row = (
            1,                          # id
            "John Doe",                 # name
            "johndoe",                  # username
            "john@example.com",         # email
            "Software developer",         # bio
            42,                         # posts
            100,                        # likes
            50,                         # following
            200,                        # followers
            True,                       # followed
            "2024-01-15 10:30:00"       # created
        )
        result = map_user(row)

        assert isinstance(result, chirp_pb2.User)
        assert result.id == 1
        assert result.name == "John Doe"
        assert result.username == "johndoe"
        assert result.email == "john@example.com"
        assert result.bio == "Software developer"
        assert result.posts == 42
        assert result.likes == 100
        assert result.following == 50
        assert result.followers == 200
        assert result.followed is True
        assert result.created == "2024-01-15 10:30:00"

    def test_map_user_with_minimal_data(self):
        row = (
            2,                          # id
            "Jane",                     # name
            "jane",                     # username
            "jane@example.com",         # email
            "",                         # bio (empty)
            0,                          # posts
            0,                          # likes
            0,                          # following
            0,                          # followers
            False,                      # followed
            "2024-01-01 00:00:00"       # created
        )
        result = map_user(row)

        assert result.id == 2
        assert result.name == "Jane"
        assert result.bio == ""
        assert result.posts == 0
        assert result.followed is False

    def test_map_user_preserves_types(self):
        row = (
            1, "Name", "username", "email@test.com",
            "Bio", 10, 20, 30, 40, True, "2024-01-01"
        )
        result = map_user(row)

        assert isinstance(result.id, int)
        assert isinstance(result.name, str)
        assert isinstance(result.username, str)
        assert isinstance(result.email, str)
        assert isinstance(result.bio, str)
        assert isinstance(result.posts, int)
        assert isinstance(result.likes, int)
        assert isinstance(result.following, int)
        assert isinstance(result.followers, int)
        assert isinstance(result.followed, bool)
        assert isinstance(result.created, str)
