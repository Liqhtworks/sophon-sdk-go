# JobProgress

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Stage** | Pointer to **string** | Current processing stage label (e.g. \&quot;probing\&quot;, \&quot;encoding\&quot;, \&quot;muxing\&quot;). | [optional] 
**Phase** | [**JobStatus**](JobStatus.md) | Canonical pipeline phase used for progress semantics. | 
**Percent** | **float32** |  | 
**PhasePercent** | Pointer to **float32** | Progress within the current phase. Null before active processing. | [optional] 
**Fps** | Pointer to **float32** |  | [optional] 
**EtaSeconds** | Pointer to **int32** |  | [optional] 
**FramesDone** | Pointer to **int32** |  | [optional] 
**FramesTotal** | Pointer to **int32** |  | [optional] 

## Methods

### NewJobProgress

`func NewJobProgress(phase JobStatus, percent float32, ) *JobProgress`

NewJobProgress instantiates a new JobProgress object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobProgressWithDefaults

`func NewJobProgressWithDefaults() *JobProgress`

NewJobProgressWithDefaults instantiates a new JobProgress object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStage

`func (o *JobProgress) GetStage() string`

GetStage returns the Stage field if non-nil, zero value otherwise.

### GetStageOk

`func (o *JobProgress) GetStageOk() (*string, bool)`

GetStageOk returns a tuple with the Stage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStage

`func (o *JobProgress) SetStage(v string)`

SetStage sets Stage field to given value.

### HasStage

`func (o *JobProgress) HasStage() bool`

HasStage returns a boolean if a field has been set.

### GetPhase

`func (o *JobProgress) GetPhase() JobStatus`

GetPhase returns the Phase field if non-nil, zero value otherwise.

### GetPhaseOk

`func (o *JobProgress) GetPhaseOk() (*JobStatus, bool)`

GetPhaseOk returns a tuple with the Phase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhase

`func (o *JobProgress) SetPhase(v JobStatus)`

SetPhase sets Phase field to given value.


### GetPercent

`func (o *JobProgress) GetPercent() float32`

GetPercent returns the Percent field if non-nil, zero value otherwise.

### GetPercentOk

`func (o *JobProgress) GetPercentOk() (*float32, bool)`

GetPercentOk returns a tuple with the Percent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPercent

`func (o *JobProgress) SetPercent(v float32)`

SetPercent sets Percent field to given value.


### GetPhasePercent

`func (o *JobProgress) GetPhasePercent() float32`

GetPhasePercent returns the PhasePercent field if non-nil, zero value otherwise.

### GetPhasePercentOk

`func (o *JobProgress) GetPhasePercentOk() (*float32, bool)`

GetPhasePercentOk returns a tuple with the PhasePercent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhasePercent

`func (o *JobProgress) SetPhasePercent(v float32)`

SetPhasePercent sets PhasePercent field to given value.

### HasPhasePercent

`func (o *JobProgress) HasPhasePercent() bool`

HasPhasePercent returns a boolean if a field has been set.

### GetFps

`func (o *JobProgress) GetFps() float32`

GetFps returns the Fps field if non-nil, zero value otherwise.

### GetFpsOk

`func (o *JobProgress) GetFpsOk() (*float32, bool)`

GetFpsOk returns a tuple with the Fps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFps

`func (o *JobProgress) SetFps(v float32)`

SetFps sets Fps field to given value.

### HasFps

`func (o *JobProgress) HasFps() bool`

HasFps returns a boolean if a field has been set.

### GetEtaSeconds

`func (o *JobProgress) GetEtaSeconds() int32`

GetEtaSeconds returns the EtaSeconds field if non-nil, zero value otherwise.

### GetEtaSecondsOk

`func (o *JobProgress) GetEtaSecondsOk() (*int32, bool)`

GetEtaSecondsOk returns a tuple with the EtaSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEtaSeconds

`func (o *JobProgress) SetEtaSeconds(v int32)`

SetEtaSeconds sets EtaSeconds field to given value.

### HasEtaSeconds

`func (o *JobProgress) HasEtaSeconds() bool`

HasEtaSeconds returns a boolean if a field has been set.

### GetFramesDone

`func (o *JobProgress) GetFramesDone() int32`

GetFramesDone returns the FramesDone field if non-nil, zero value otherwise.

### GetFramesDoneOk

`func (o *JobProgress) GetFramesDoneOk() (*int32, bool)`

GetFramesDoneOk returns a tuple with the FramesDone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFramesDone

`func (o *JobProgress) SetFramesDone(v int32)`

SetFramesDone sets FramesDone field to given value.

### HasFramesDone

`func (o *JobProgress) HasFramesDone() bool`

HasFramesDone returns a boolean if a field has been set.

### GetFramesTotal

`func (o *JobProgress) GetFramesTotal() int32`

GetFramesTotal returns the FramesTotal field if non-nil, zero value otherwise.

### GetFramesTotalOk

`func (o *JobProgress) GetFramesTotalOk() (*int32, bool)`

GetFramesTotalOk returns a tuple with the FramesTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFramesTotal

`func (o *JobProgress) SetFramesTotal(v int32)`

SetFramesTotal sets FramesTotal field to given value.

### HasFramesTotal

`func (o *JobProgress) HasFramesTotal() bool`

HasFramesTotal returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


