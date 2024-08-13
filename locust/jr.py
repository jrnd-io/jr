"""
JR Locust library
"""
import time
import subprocess
from locust import User


class JRUser(User):
    """
    Superclass for JR users
    """
    abstract = True

    def __init__(self, environment):
        """
        Init User
        """
        super().__init__(environment)
        self.env = environment
        self.jr = "/usr/bin/jr"

    def run_jr(self,  args):
        """
        Runs JR
        """
        start_perf_counter = time.perf_counter()
        exc = None
        try:
            subprocess.run([self.jr]+args, check=True)
        except subprocess.CalledProcessError as e:
            exc = e

        self.env.events.request.fire(
                request_type="jr",
                name="jr",
                response_time=(time.perf_counter() - start_perf_counter) * 1000,
                response_length=0,
                response=None,
                context=None,
                exception=exc)
