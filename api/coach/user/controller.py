from pydantic import BaseModel, Field
from fastapi.errors import ValidationError
from fastapi.responses import JSONResponse
from fastapi import APIRouter


router = APIRouter(prefix="/users")

class UserCreationRequest(BaseModel):
    """A specification of the user to be created"""
    name: str = field(min=1, max=50, pattern="^[a-zA-Z- ]{1,50}$")


class ErrorResponse(BaseModel):
    error: str
    description: str

error_resp_by_validation_error = {
    "invalid length": ("invalid_length", "Field '{}' is too short")
}

@router.exception_handler(RequestValidationError)
def handle_validation_errors(request, exc):
    error_type = exc.error()[0]['type']
    error = error_response_by_type.get(error_type)
    if error is not None:
        error_code, error_description_template = error_response_by_type[error_type]
        error_description = error_description_template.format()
        error_response = ErrorResponse(error=error_code, description=error_description)
    else:
        error_response = ErrorResponse(error="client_error", description="a client error occurred")
    return JSONResponse(status_code=422, content=jsonable_encoder(error_response))


@router.post("/", status_code=201)
def create_user(self, request: UserCreationRequest) -> None:
    """Creates a new user according to the specified request.

    Args:
        request (UserCreationRequest): A specification of the user to be created

    Raises:
        InvalidNameError: If the specified name does not match name restrictions
    """
    try:
        created_user = user_service.save(user)
    except Exception as e:
        log.error("An error occurred")
        raise GeneralError()
        
    return created_user 
