<oap:target id="user-management-module">
  <oap:target id="sqlalchemy-models-section">
    <oap:target id="user-model">
      <oap:target id="user-model-name">User</oap:target>
      <oap:target id="user-id-field">id</oap:target>
      <oap:target id="user-email-field">email</oap:target>
      <oap:target id="user-name-field">name</oap:target>
      <oap:target id="user-hashed-password-field">hashed_password</oap:target>
      <oap:target id="user-role-field">role</oap:target>
      <oap:target id="user-is-active-field">is_active</oap:target>
      <oap:target id="user-created-at-field">created_at</oap:target>
      from datetime import datetime

      from sqlalchemy import Boolean, Column, DateTime, Integer, String

      from .database import Base


      class User(Base):
          __tablename__ = "users"

          id = Column(Integer, primary_key=True, index=True)
          email = Column(String, unique=True, index=True, nullable=False)
          name = Column(String, nullable=False)
          hashed_password = Column(String, nullable=False)
          role = Column(String, nullable=False, default="user")
          is_active = Column(Boolean, default=True, nullable=False)
          created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    </oap:target>
  </oap:target>

  <oap:target id="pydantic-schemas-section">
    <oap:target id="user-create-schema">from typing import Literal


      <oap:target id="user-create-schema-name">UserCreate</oap:target>
      <oap:target id="user-create-email-field">email</oap:target>
      <oap:target id="user-create-name-field">name</oap:target>
      <oap:target id="user-create-password-field">password</oap:target>
      <oap:target id="user-create-role-field">          role: Literal["admin", "editor", "viewer"] = "viewer"</oap:target>
      <oap:target id="user-create-is-active-field">is_active</oap:target>
      from pydantic import BaseModel, ConfigDict, EmailStr, Field


      class UserCreate(BaseModel):
          email: EmailStr
          name: str = Field(min_length=1, max_length=255)
          password: str = Field(min_length=8, max_length=255)
          role: str = Field(default="user", max_length=50)
          is_active: bool = True
    </oap:target>

    <oap:target id="user-update-schema">
      <oap:target id="user-update-schema-name">UserUpdate</oap:target>
      <oap:target id="user-update-email-field">email</oap:target>
      <oap:target id="user-update-name-field">name</oap:target>
      <oap:target id="user-update-password-field">password</oap:target>
      <oap:target id="user-update-role-field">role</oap:target>
      <oap:target id="user-update-is-active-field">is_active</oap:target>
      class UserUpdate(BaseModel):
          email: EmailStr | None = None
          name: str | None = Field(default=None, min_length=1, max_length=255)
          password: str | None = Field(default=None, min_length=8, max_length=255)
          role: str | None = Field(default=None, max_length=50)
          is_active: bool | None = None
    </oap:target>

    <oap:target id="user-response-schema">
      <oap:target id="user-response-schema-name">UserResponse</oap:target>
      <oap:target id="user-response-id-field">id</oap:target>
      <oap:target id="user-response-email-field">email</oap:target>
      <oap:target id="user-response-name-field">name</oap:target>
      <oap:target id="user-response-role-field">role</oap:target>
      <oap:target id="user-response-is-active-field">is_active</oap:target>
      <oap:target id="user-response-created-at-field">created_at</oap:target>
      class UserResponse(BaseModel):
          model_config = ConfigDict(from_attributes=True)

          id: int
          email: EmailStr
          name: str
          role: str
          is_active: bool
          created_at: datetime
    </oap:target>

    <oap:target id="user-list-schema">
      <oap:target id="user-list-schema-name">UserList</oap:target>
      <oap:target id="user-list-items-field">items</oap:target>
      <oap:target id="user-list-total-field">total</oap:target>
      <oap:target id="user-list-page-field">page</oap:target>
      <oap:target id="user-list-size-field">size</oap:target>
      class UserList(BaseModel):
          items: list[UserResponse]
          total: int
          page: int
          size: int
    </oap:target>
  </oap:target>

  <oap:target id="crud-functions-section">
    <oap:target id="create-user-function">
      <oap:target id="create-user-function-name">create_user</oap:target>
      <oap:target id="create-user-email-check">email uniqueness check</oap:target>
      <oap:target id="create-user-password-hashing">password hashing</oap:target>
      <oap:target id="create-user-error">HTTPException</oap:target>
      from fastapi import HTTPException, status
      from sqlalchemy.orm import Session

      from .models import User
      from .schemas import UserCreate, UserUpdate


      def create_user(db: Session, user_in: UserCreate) -> User:
          existing_user = db.query(User).filter(User.email == user_in.email).first()
          if existing_user:
              raise HTTPException(
                  status_code=status.HTTP_400_BAD_REQUEST,
                  detail="Email already registered",
              )

          user = User(
              email=user_in.email,
              name=user_in.name,
              hashed_password=hash_password(user_in.password),
              role=user_in.role,
              is_active=user_in.is_active,
          )
          db.add(user)
          db.commit()
          db.refresh(user)
          return user
    </oap:target>

    <oap:target id="get-user-function">
      <oap:target id="get-user-function-name">get_user</oap:target>
      <oap:target id="get-user-not-found-error">HTTPException</oap:target>
      def get_user(db: Session, user_id: int) -> User:
          user = db.query(User).filter(User.id == user_id).first()
          if not user:
              raise HTTPException(
                  status_code=status.HTTP_404_NOT_FOUND,
                  detail="User not found",
              )
          return user
    </oap:target>

    <oap:target id="list-users-function">def list_users(db: Session, page: int = 1, size: int = 10, role: str | None = None, is_active: bool | None = None) -> tuple[list[User], int]:
    offset = (page - 1) * size
    query = db.query(User)

    if role is not None:
        query = query.filter(User.role == role)
    if is_active is not None:
        query = query.filter(User.is_active == is_active)

    total = query.count()
    users = query.order_by(User.id).offset(offset).limit(size).all()
    return users, total</oap:target>

    <oap:target id="update-user-function">
      <oap:target id="update-user-function-name">update_user</oap:target>
      <oap:target id="update-user-email-check">email uniqueness check</oap:target>
      <oap:target id="update-user-error">HTTPException</oap:target>
      def update_user(db: Session, user_id: int, user_in: UserUpdate) -> User:
          user = get_user(db, user_id)

          if user_in.email and user_in.email != user.email:
              existing_user = db.query(User).filter(User.email == user_in.email).first()
              if existing_user:
                  raise HTTPException(
                      status_code=status.HTTP_400_BAD_REQUEST,
                      detail="Email already registered",
                  )
              user.email = user_in.email

          if user_in.name is not None:
              user.name = user_in.name
          if user_in.password is not None:
              user.hashed_password = hash_password(user_in.password)
          if user_in.role is not None:
              user.role = user_in.role
          if user_in.is_active is not None:
              user.is_active = user_in.is_active

          db.commit()
          db.refresh(user)
          return user
    </oap:target>

    <oap:target id="delete-user-function">
      <oap:target id="delete-user-function-name">delete_user</oap:target>
      <oap:target id="delete-user-error">HTTPException</oap:target>
      def delete_user(db: Session, user_id: int) -> None:
          user = get_user(db, user_id)
          db.delete(user)
          db.commit()
    </oap:target>

    <oap:target id="password-helper-function">
      <oap:target id="password-helper-name">hash_password</oap:target>
      <oap:target id="password-helper-placeholder">placeholder hashing implementation</oap:target>
      def hash_password(password: str) -> str:
          return password
    </oap:target>
  </oap:target>

  <oap:target id="fastapi-router-section">
    <oap:target id="router-prefix">/users</oap:target>
    <oap:target id="router-tags">users</oap:target>
    <oap:target id="create-user-endpoint">
      <oap:target id="create-user-method">POST</oap:target>
      <oap:target id="create-user-path">/users</oap:target>
      @router.post("/", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
      def create_user_endpoint(user_in: UserCreate, db: Session = Depends(get_db)):
          return create_user(db, user_in)
    </oap:target>

    <oap:target id="list-users-endpoint">
      <oap:target id="list-users-method">GET</oap:target>
      <oap:target id="list-users-path">/users</oap:target>
      @router.get("/", response_model=UserList)
      def list_users_endpoint(
          page: int = Query(default=1, ge=1),
          size: int = Query(default=10, ge=1, le=100),
          db: Session = Depends(get_db),
      ):
          users, total = list_users(db, page=page, size=size)
          return UserList(items=users, total=total, page=page, size=size)
    </oap:target>

    <oap:target id="get-user-endpoint">
      <oap:target id="get-user-method">GET</oap:target>
      <oap:target id="get-user-path">/users/{id}</oap:target>
      @router.get("/{user_id}", response_model=UserResponse)
      def get_user_endpoint(user_id: int, db: Session = Depends(get_db)):
          return get_user(db, user_id)
    </oap:target>

    <oap:target id="update-user-endpoint">
      <oap:target id="update-user-method">PUT</oap:target>
      <oap:target id="update-user-path">/users/{id}</oap:target>
      @router.put("/{user_id}", response_model=UserResponse)
      def update_user_endpoint(user_id: int, user_in: UserUpdate, db: Session = Depends(get_db)):
          return update_user(db, user_id, user_in)
    </oap:target>

    <oap:target id="delete-user-endpoint">
      <oap:target id="delete-user-method">DELETE</oap:target>
      <oap:target id="delete-user-path">/users/{id}</oap:target>
      @router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
      def delete_user_endpoint(user_id: int, db: Session = Depends(get_db)):
          delete_user(db, user_id)
          return None
    </oap:target>

    <oap:target id="router-code">
      from fastapi import APIRouter, Depends, Query, status
      from sqlalchemy.orm import Session

      from .crud import create_user, delete_user, get_user, list_users, update_user
      from .database import get_db
      from .schemas import UserCreate, UserList, UserResponse, UserUpdate

      router = APIRouter(prefix="/users", tags=["users"])
    </oap:target>
  </oap:target>

  <oap:target id="error-handling-section">
    <oap:target id="http-exception-usage">HTTPException</oap:target>
    <oap:target id="not-found-status-code">404</oap:target>
    <oap:target id="conflict-status-code">400</oap:target>
    <oap:target id="validation-error-handling">Pydantic validation</oap:target>
    <oap:target id="crud-error-contract">raise HTTPException for not found and duplicate email cases</oap:target>
  </oap:target>
</oap:target>