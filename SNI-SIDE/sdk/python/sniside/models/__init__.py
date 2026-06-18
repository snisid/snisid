class BaseService:
    def __init__(self, http):
        self.http = http
    
    def get(self, path, **kwargs):
        return self.http.get(path, **kwargs)
    
    def post(self, path, data=None):
        return self.http.post(path, data)
    
    def put(self, path, data=None):
        return self.http.put(path, data)
    
    def delete(self, path):
        return self.http.delete(path)
