// Code generated by go-swagger; DO NOT EDIT.

package daemon

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetDebuginfoParams creates a new GetDebuginfoParams object
// with the default values initialized.
func NewGetDebuginfoParams() *GetDebuginfoParams {

	return &GetDebuginfoParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetDebuginfoParamsWithTimeout creates a new GetDebuginfoParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetDebuginfoParamsWithTimeout(timeout time.Duration) *GetDebuginfoParams {

	return &GetDebuginfoParams{

		timeout: timeout,
	}
}

// NewGetDebuginfoParamsWithContext creates a new GetDebuginfoParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetDebuginfoParamsWithContext(ctx context.Context) *GetDebuginfoParams {

	return &GetDebuginfoParams{

		Context: ctx,
	}
}

// NewGetDebuginfoParamsWithHTTPClient creates a new GetDebuginfoParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetDebuginfoParamsWithHTTPClient(client *http.Client) *GetDebuginfoParams {

	return &GetDebuginfoParams{
		HTTPClient: client,
	}
}

/*GetDebuginfoParams contains all the parameters to send to the API endpoint
for the get debuginfo operation typically these are written to a http.Request
*/
type GetDebuginfoParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get debuginfo params
func (o *GetDebuginfoParams) WithTimeout(timeout time.Duration) *GetDebuginfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get debuginfo params
func (o *GetDebuginfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get debuginfo params
func (o *GetDebuginfoParams) WithContext(ctx context.Context) *GetDebuginfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get debuginfo params
func (o *GetDebuginfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get debuginfo params
func (o *GetDebuginfoParams) WithHTTPClient(client *http.Client) *GetDebuginfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get debuginfo params
func (o *GetDebuginfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetDebuginfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
