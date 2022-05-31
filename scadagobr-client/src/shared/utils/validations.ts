import type {ValidationError} from "yup";

export function extractErrors(err: ValidationError): Record<string, string> {
    return err.inner.reduce((acc, err) => {
        return {...acc, [err.path]: err.message};
    }, {});
}

type Errors = {
    [key: string]: string;
};

export async function getErrorsFromSchema(schema, data): Promise<undefined | Errors> {

    try {
        await schema.validateSync(data, {abortEarly: false});
    } catch (err) {
        return extractErrors(err)
    }

    return null
}

export function inputError(error) : any {
    const response = {
        type: "",
        message: ""
    }

    if (error){
        response.type = "is-danger"
        response.message = error
    }

    return response;
}