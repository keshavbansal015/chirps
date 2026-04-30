from authservice.mappers import map_session
from authservice import chirp_pb2


class TestMapSession:
    def test_map_session_with_valid_data(self):
        row = ('session-123', 42, '2024-01-15 10:30:00')
        result = map_session(row)

        assert isinstance(result, chirp_pb2.Session)
        assert result.id == 'session-123'
        assert result.user_id == 42
        assert result.created == '2024-01-15 10:30:00'

    def test_map_session_with_different_data(self):
        row = ('another-session', 999, '2023-12-25 00:00:00')
        result = map_session(row)

        assert result.id == 'another-session'
        assert result.user_id == 999
        assert result.created == '2023-12-25 00:00:00'

    def test_map_session_preserves_types(self):
        row = ('test-id', 1, '2024-01-01')
        result = map_session(row)

        assert isinstance(result.id, str)
        assert isinstance(result.user_id, int)
        assert isinstance(result.created, str)
