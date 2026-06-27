from sqlalchemy import Column, Integer, String, Float, DateTime, Enum
from sqlalchemy.ext.declarative import declarative_base
import datetime
from app.utils.errors import JobsBotError

Base = declarative_base()


class Run(Base):
    __tablename__ = "runs"
    id = Column(String, primary_key=True)
    state = Column(String)
    started_at = Column(DateTime, default=datetime.datetime.utcnow)
    completed_at = Column(DateTime, nullable=True)
    max_applications = Column(Integer)
    platforms = Column(String)

    def save(self, session):
        try:
            if not self.id or not isinstance(self.id, str):
                raise JobsBotError("VALID901", {"run_id": self.id})
            if not self.state or not isinstance(self.state, str):
                raise JobsBotError("VALID902", {"state": self.state})
            session.add(self)
            session.commit()
        except Exception as e:
            raise JobsBotError("DB002", {"run_id": self.id, "error": str(e)})


class Application(Base):
    __tablename__ = "applications"
    id = Column(String, primary_key=True)
    run_id = Column(String)
    job_url = Column(String)
    state = Column(String)
    score = Column(Float)
    reason = Column(String)
    applied_at = Column(DateTime, nullable=True)

    def save(self, session):
        try:
            if not self.id or not isinstance(self.id, str):
                raise JobsBotError("VALID903", {"application_id": self.id})
            if not self.run_id or not isinstance(self.run_id, str):
                raise JobsBotError("VALID904", {"run_id": self.run_id})
            if not self.job_url or not isinstance(self.job_url, str):
                raise JobsBotError("VALID905", {"job_url": self.job_url})
            session.add(self)
            session.commit()
        except Exception as e:
            raise JobsBotError(
                "DB002", {"application_id": self.id, "error": str(e)})


class Log(Base):
    __tablename__ = "logs"
    id = Column(Integer, primary_key=True, autoincrement=True)
    run_id = Column(String)
    event = Column(String)
    timestamp = Column(DateTime, default=datetime.datetime.utcnow)
    details = Column(String)

    def save(self, session):
        try:
            if not self.run_id or not isinstance(self.run_id, str):
                raise JobsBotError("VALID906", {"run_id": self.run_id})
            if not self.event or not isinstance(self.event, str):
                raise JobsBotError("VALID907", {"event": self.event})
            session.add(self)
            session.commit()
        except Exception as e:
            raise JobsBotError("DB002", {"log_id": self.id, "error": str(e)})
