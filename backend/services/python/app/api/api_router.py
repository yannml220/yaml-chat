from fastapi import APIRouter
from .routers import nlp_router

api_router = APIRouter()

api_router.include_router(nlp_router, prefix="/nlp", tags=["nlp"])
