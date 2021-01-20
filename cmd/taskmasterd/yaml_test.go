package main

import "testing"

func strToPointer(str string) *string {
	return &str
}

func intToPointer(nb int) *int {
	return &nb
}

func boolToPointer(b bool) *bool {
	return &b
}

func autorestartTypeToPointer(t AutorestartType) *AutorestartType {
	return &t
}

func stopSignalToPointer(signal StopSignal) *StopSignal {
	return &signal
}

func TestCmdIsRequired(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd: nil,
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Cmd" && validationError.Issue == ValidationIssueEmptyField) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Cmd",
				ValidationIssueEmptyField,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestNumprocsSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:      strToPointer("cmd"),
				Numprocs: nil,
			},
		},
	}

	config, _ := programs.Validate()

	if numprocs := config["taskmaster"].Numprocs; numprocs != 1 {
		t.Errorf(
			"Numprocs not set to correct default value: %v; expected %v",
			numprocs,
			1,
		)
	}
}

func TestNumprocsIsNotOutsideLowerBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:      strToPointer("cmd"),
				Numprocs: intToPointer(-1),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Numprocs" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Numprocs",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestNumprocsIsNotOutsideUpperBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:      strToPointer("cmd"),
				Numprocs: intToPointer(200),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Numprocs" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Numprocs",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestAutostartSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:       strToPointer("cmd"),
				Numprocs:  intToPointer(10),
				Autostart: nil,
			},
		},
	}

	config, _ := programs.Validate()

	if autostart := config["taskmaster"].Autostart; autostart != true {
		t.Errorf(
			"Autostart not set to correct default value: %v; expected %v",
			autostart,
			1,
		)
	}
}

func TestAutorestartSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:         strToPointer("cmd"),
				Numprocs:    intToPointer(10),
				Autostart:   boolToPointer(true),
				Autorestart: nil,
			},
		},
	}

	config, _ := programs.Validate()

	if autorestart := config["taskmaster"].Autorestart; autorestart != AutorestartUnexpected {
		t.Errorf(
			"Autorestart not set to correct default value: %v; expected %v",
			autorestart,
			AutorestartUnexpected,
		)
	}
}

func TestAutorestartIsValidValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:         strToPointer("cmd"),
				Numprocs:    intToPointer(10),
				Autostart:   boolToPointer(true),
				Autorestart: autorestartTypeToPointer(AutorestartOn),
			},
		},
	}

	_, err := programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for autorestart = %s", AutorestartOn)
	}

	*programs.Programs["taskmaster"].Autorestart = AutorestartOff

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for autorestart = %s", AutorestartOff)
	}

	*programs.Programs["taskmaster"].Autorestart = AutorestartUnexpected

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for autorestart = %s", AutorestartUnexpected)
	}

	*programs.Programs["taskmaster"].Autorestart = "Invalid value"

	_, err = programs.Validate()
	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Autorestart" && validationError.Issue == ValidationIssueUnexpectedValue) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Autorestart",
				ValidationIssueUnexpectedValue,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStarttimeSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:         strToPointer("cmd"),
				Numprocs:    intToPointer(10),
				Autostart:   boolToPointer(true),
				Autorestart: autorestartTypeToPointer(AutorestartOn),
				Starttime:   nil,
			},
		},
	}

	config, _ := programs.Validate()

	if starttime := config["taskmaster"].Starttime; starttime != 1 {
		t.Errorf(
			"Starttime not set to correct default value: %v; expected %v",
			starttime,
			1,
		)
	}
}

func TestStarttimeIsNotOutsideLowerBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:         strToPointer("cmd"),
				Numprocs:    intToPointer(10),
				Autostart:   boolToPointer(true),
				Autorestart: autorestartTypeToPointer(AutorestartOn),
				Starttime:   intToPointer(-1),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Starttime" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Starttime",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStarttimeIsNotOutsideUpperBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:         strToPointer("cmd"),
				Numprocs:    intToPointer(10),
				Autostart:   boolToPointer(true),
				Autorestart: autorestartTypeToPointer(AutorestartOn),
				Starttime:   intToPointer(100000),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Starttime" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Starttime",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStartretriesSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: nil,
			},
		},
	}

	config, _ := programs.Validate()

	if startretries := config["taskmaster"].Startretries; startretries != 3 {
		t.Errorf(
			"Startretries not set to correct default value: %v; expected %v",
			startretries,
			3,
		)
	}
}

func TestStartretriesIsNotOutsideLowerBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(-1),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Startretries" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Startretries",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStartretriesIsNotOutsideUpperBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(50),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Startretries" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Startretries",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStopsignalSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   nil,
			},
		},
	}

	config, _ := programs.Validate()

	if stopsignal := config["taskmaster"].Stopsignal; stopsignal != StopSignalTerm {
		t.Errorf(
			"Stopsignal not set to correct default value: %v; expected %v",
			stopsignal,
			StopSignalTerm,
		)
	}
}

func TestStopsignalIsValidValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
			},
		},
	}

	_, err := programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalTerm)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalHup

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalHup)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalInt

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalInt)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalQuit

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalQuit)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalKill

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalKill)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalUsr1

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalUsr1)
	}

	*programs.Programs["taskmaster"].Stopsignal = StopSignalUsr2

	_, err = programs.Validate()
	if err != nil {
		t.Errorf("Expected no error for Stopsignal = %s", StopSignalUsr2)
	}

	*programs.Programs["taskmaster"].Stopsignal = "Invalid value"

	_, err = programs.Validate()
	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Stopsignal" && validationError.Issue == ValidationIssueUnexpectedValue) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Stopsignal",
				ValidationIssueUnexpectedValue,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStoptimeSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     nil,
			},
		},
	}

	config, _ := programs.Validate()

	if stoptime := config["taskmaster"].Stoptime; stoptime != 10 {
		t.Errorf(
			"Stoptime not set to correct default value: %v; expected %v",
			stoptime,
			intToPointer(10),
		)
	}
}

func TestStoptimeIsNotOutsideLowerBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     intToPointer(-1),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Stoptime" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Stoptime",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStoptimeIsNotOutsideUpperBounds(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     intToPointer(5000),
			},
		},
	}

	_, err := programs.Validate()
	if err == nil {
		t.Errorf("Validate should have returned an error")
	}

	if validationError, ok := err.(*ErrProgramsYamlValidation); ok {
		if !(validationError.Field == "Programs[taskmaster].Stoptime" && validationError.Issue == ValidationIssueValueOutsideBounds) {
			t.Errorf(
				"Incorrect error: (%s, %s); expected (%s, %s)",
				validationError.Field,
				validationError.Issue,
				"Programs[taskmaster].Stoptime",
				ValidationIssueValueOutsideBounds,
			)
			return
		}
		return
	}

	t.Errorf("Returned invalid error")
}

func TestStdoutSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     intToPointer(60),
				Stdout:       nil,
			},
		},
	}

	config, _ := programs.Validate()

	if stdout := config["taskmaster"].Stdout; stdout != string(StdTypeAuto) {
		t.Errorf(
			"Stdout not set to correct default value: %v; expected %v",
			stdout,
			string(StdTypeAuto),
		)
	}
}

func TestStderrSetToDefaultValue(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("cmd"),
				Numprocs:     intToPointer(10),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Starttime:    intToPointer(5),
				Startretries: intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     intToPointer(60),
				Stderr:       nil,
			},
		},
	}

	config, _ := programs.Validate()

	if stderr := config["taskmaster"].Stderr; stderr != string(StdTypeAuto) {
		t.Errorf(
			"Stderr not set to correct default value: %v; expected %v",
			stderr,
			string(StdTypeAuto),
		)
	}
}

func TestParsesValidFullConfiguration(t *testing.T) {
	programs := ProgramsYaml{
		Programs: map[string]ProgramYaml{
			"taskmaster": {
				Cmd:          strToPointer("echo"),
				Numprocs:     intToPointer(1),
				Umask:        strToPointer("066"),
				Workingdir:   strToPointer("/dir"),
				Autostart:    boolToPointer(true),
				Autorestart:  autorestartTypeToPointer(AutorestartOn),
				Exitcodes:    []int{0},
				Startretries: intToPointer(3),
				Starttime:    intToPointer(10),
				Stopsignal:   stopSignalToPointer(StopSignalTerm),
				Stoptime:     intToPointer(10),
				Stdout:       strToPointer("/dev/stdout"),
				Stderr:       strToPointer("/dev/stderr"),
				Env: map[string]string{
					"TERM": "DUMB",
				},
			},
		},
	}

	_, err := programs.Validate()
	if err != nil {
		t.Errorf("Validation error on valid configuration: %v", err)
	}
}
