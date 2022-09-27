package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"github.com/thecathe/gocurrency_tool/analyser/scopemanager"
)

var scope_manager *scopemanager.ScopeManager

// called on every file
func AnalyseAst(fileset *token.FileSet, package_name string, filename string, node ast.Node, channel chan Counter, name string) {

	var counter Counter = Counter{Go_count: 0, Send_count: 0, Rcv_count: 0, Chan_count: 0, filename: name}

	if _scope_manager, ok := scopemanager.NewScopeManager(filename, fileset); ok == nil {
		scope_manager = _scope_manager
		DebugLog("Analyser, %s: Scope Manager Created.\n", filename)
	} else {
		FailureLog("Analyser, %s: Scope Manager Creation Failed\n", filename)
		return
	}

	GeneralLog("Analyser, %s: Printing ScopeManager.Scopes...\n%s\n", filename, scope_manager.ScopeMap.ToString())

	// go through file
	switch file := node.(type) {
	case *ast.File:
		// for each file

		if _scope_manager, _parse_type := scope_manager.ParseNode(&node); _parse_type != scopemanager.PARSE_NONE {
			scope_manager = _scope_manager
			DebugLog("Analyser, %s: File Scope Successful.\n", filename)

			// generate scope map
			for _, file_decl := range file.Decls{
				// go through each scope in file
				if _scope_manager, _parse_type := scope_manager.ParseNode(&file_decl); _parse_type  != scopemanager.PARSE_NONE {
					scope_manager = _scope_manager
					DebugLog("Analyser, %s: Global Scope Successful.\n", filename)
					var node_visit_count int = 0

					//
					// start of each node inspect
					//
					ast.Inspect(file_decl, func(_decl ast.Node) bool {
						node_visit_count++
						// for every node in this scope
						// give current node to manager
						if _scope_manager, parse_type := scope_manager.ParseNode(&_decl); parse_type != scopemanager.PARSE_NONE {
							scope_manager = _scope_manager
							DebugLog("Analyser, %s: Inspect: %d, Type %s Success.\n%s\n", filename, node_visit_count, parse_type)

						} else {
							scope_manager = _scope_manager

							DebugLog("Analyser, %s: Inspect: %d, Type %s Failed.\n", filename, node_visit_count, parse_type)
						}

						return true
					})
					//
					// end of each node insepct
					//

				} else {
					FailureLog("Analyser, %s: Global Scope Unsuccessful: %s\n", filename, _parse_type)
				}
			}
		} else {
			FailureLog("Analyser, %s: File Scope Unsuccessful: %s\n", filename, _parse_type)
		}

		// // add scope
		// if _scope_manager, scope_id, ok := scope_manager.NewScope(file, *scopemanager.SCOPE_TYPE_FILE); ok == nil {
		// 	scope_manager = _scope_manager
		// 	DebugLog("Analyser, %s: File Scope Successful.\n", filename)

		// 	// generate scope map
		// 	for _, file_decl := range file.Decls {
		// 		// for each global decl: func, const, var, import
		// 		if _scope_manager, scope_id, ok := scope_manager.NewScope(file_decl); ok == nil {
		// 			scope_manager = _scope_manager
		// 			DebugLog("\nAnalyser, %s: Global Decl Scope Successful.\n", filename)
		// 			var node_visit_count int = 0
		// 			//
		// 			// start of each node inspect
		// 			//
		// 			ast.Inspect(file_decl, func(_decl ast.Node) bool {
		// 				node_visit_count++
		// 				// for every node in this scope
		// 				// give current node to manager
		// 				if _scope_manager, parse_type, type_id, ok := scope_manager.ParseNode(_decl); ok == nil {
		// 					scope_manager = _scope_manager

		// 					var parse_type_string string = type_id
		// 					switch parse_type {
		// 					case *scopemanager.PARSE_NONE:
		// 						// error occured
		// 					case *scopemanager.PARSE_SCOPE:
		// 						// scope found
		// 						// parse_type_string = string(scope_manager.Scopes[ScopeID(type_id)].ID)
		// 					case *scopemanager.PARSE_DECL:
		// 						// decl found
		// 					case *scopemanager.PARSE_ASSIGN:
		// 						// assignment found
		// 					}
		// 					DebugLog("Analyser, %s: Inspect: %d, Type %s Success.\n%s\n", filename, node_visit_count, parse_type, parse_type_string)

		// 				} else {
		// 					scope_manager = _scope_manager

		// 					DebugLog("Analyser, %s: Inspect: %d, Type %s Failed.\n", filename, node_visit_count, parse_type)
		// 				}

		// 				return true
		// 			})
		// 			//
		// 			// end of each node insepct
		// 			//
		// 		} else {
		// 			scope_manager = _scope_manager
		// 			FailureLog("Analyser, %s: Decl Setup, NewScope Failed...\n\tID: %s\n", filename, scope_id)
		// 			continue
		// 		}
		// 		// finished that global decl
		// 		if _scope_manager, scope_id, ok := scope_manager.PopStack(); ok == nil {
		// 			scope_manager = _scope_manager
		// 		} else {
		// 			scope_manager = _scope_manager
		// 			FailureLog("Analyser, %s: Decl Pop, PopStack Failed...\n\tID: %s\n", filename, scope_id)
		// 		}
		// 	}

		// } else {
		// 	FailureLog("Analyser, %s: File Scope Unsuccessful.\n", filename)
		// }

		// if _scope_manager, scope_id, ok := scope_manager.NewScope(file); ok == nil {
		// 	scope_manager = _scope_manager
		// 	DebugLog("Analyser, %s: File Scope Successful.\n", filename)
		// 	// generate scope map
		// 	for _, file_decl := range file.Decls {
		// 		// for each global decl: func, const, var, import
		// 		if _scope_manager, scope_id, ok := scope_manager.NewScope(file_decl); ok == nil {
		// 			scope_manager = _scope_manager
		// 			DebugLog("\nAnalyser, %s: Global Decl Scope Successful.\n", filename)
		// 			var node_visit_count int = 0
		// 			//
		// 			// start of each node inspect
		// 			//
		// 			ast.Inspect(file_decl, func(_decl ast.Node) bool {
		// 				node_visit_count++
		// 				// for every node in this scope
		// 				// give current node to manager
		// 				if _scope_manager, parse_type, type_id, ok := scope_manager.ParseNode(_decl); ok == nil {
		// 					scope_manager = _scope_manager

		// 					var parse_type_string string = type_id
		// 					switch parse_type {
		// 					case *scopemanager.PARSE_NONE:
		// 						// error occured
		// 					case *scopemanager.PARSE_SCOPE:
		// 						// scope found
		// 						// parse_type_string = string(scope_manager.Scopes[ScopeID(type_id)].ID)
		// 					case *scopemanager.PARSE_DECL:
		// 						// decl found
		// 					case *scopemanager.PARSE_ASSIGN:
		// 						// assignment found
		// 					}
		// 					DebugLog("Analyser, %s: Inspect: %d, Type %s Success.\n%s\n", filename, node_visit_count, parse_type, parse_type_string)

		// 				} else {
		// 					scope_manager = _scope_manager

		// 					DebugLog("Analyser, %s: Inspect: %d, Type %s Failed.\n", filename, node_visit_count, parse_type)
		// 				}

		// 				return true
		// 			})
		// 			//
		// 			// end of each node insepct
		// 			//
		// 		} else {
		// 			scope_manager = _scope_manager
		// 			FailureLog("Analyser, %s: Decl Setup, NewScope Failed...\n\tID: %s\n", filename, scope_id)
		// 			continue
		// 		}
		// 		// finished that global decl
		// 		if _scope_manager, scope_id, ok := scope_manager.PopStack(); ok == nil {
		// 			scope_manager = _scope_manager
		// 		} else {
		// 			scope_manager = _scope_manager
		// 			FailureLog("Analyser, %s: Decl Pop, PopStack Failed...\n\tID: %s\n", filename, scope_id)
		// 		}
		// 	}
		// } else {
		// 	scope_manager = _scope_manager
		// 	FailureLog("Analyser, %s: File Setup, NewScope Failed...\n\tID: %s\n", filename, scope_id)
		// }
	}

	var env []string = []string{}

	// then analyse each node
	switch file := node.(type) {
	case *ast.File:
		addGlobalVarToEnv(file, &env)
		for _, decl := range file.Decls {
			fresh_env := env
			ast.Inspect(decl, func(decl ast.Node) bool {
				analyseNode(fileset, package_name, filename, decl, &counter, &fresh_env)
				return true
			})
		}
	}

	GeneralLog("Analyser, %s: Finished Analysis.\n", filename)
	setFeaturesNumber(&counter)
	GeneralLog("Analyser, %s: Finished, Returning.\n", filename)
	channel <- counter
}

func analyseNode(fileset *token.FileSet, package_name string, filename string, node ast.Node, counter *Counter, env *[]string) {

	var feature Feature = Feature{
		F_filename:     filename,
		F_package_name: package_name,
		F_type:         NONE}

	switch x := node.(type) {
	// if generic declaration that usese ( )
	case *ast.GenDecl:
		// if variable declaration
		if x.Tok == token.VAR {
			for _, spec := range x.Specs {
				switch value_spec := spec.(type) {
				// if it is either a constant or variable
				case *ast.ValueSpec:
					for index, value := range value_spec.Values {
						switch call_expr := value.(type) {
						// if it has arguments
						case *ast.CallExpr:
							switch ident := call_expr.Fun.(type) {
							// get its identifier
							case *ast.Ident:
								// if it is a channel declaration
								if ident.Name == "make" {
									// if it is a valid declaration
									if len(call_expr.Args) > 0 {
										switch call_expr.Args[0].(type) {
										// if it is a channel
										case *ast.ChanType:
											ident1 := value_spec.Names[index]
											*env = append(*env, ident1.Name)
											checkDepthChan(call_expr, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
										}
									}
								}
							}
						// if it has elements for arguments
						case *ast.CompositeLit:
							switch array_type := call_expr.Type.(type) {
							// get its identifier
							case *ast.Ident:
								// Possible assignment of a struct struct = Struct{bla:0, bla1}
								// for each element,
								for _, elt := range call_expr.Elts {
									switch valueExp := elt.(type) {
									// with their keys,
									case *ast.KeyValueExpr:
										switch ident := valueExp.Key.(type) {
										// get their keys identifier
										case *ast.Ident:
											switch call := valueExp.Value.(type) {
											// if the element has arguments
											case *ast.CallExpr:
												ident1 := value_spec.Names[index]
												checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
											}
										}
									}
								}

							case *ast.ArrayType:
								checkArrayType(array_type, counter, feature, fileset, 1)
							case *ast.MapType:
								chan_in_map := false
								// we have a declaration of a map
								switch array_type.Key.(type) {
								case *ast.ChanType:
									chan_in_map = true
								}
								switch array_type.Value.(type) {
								case *ast.ChanType:
									chan_in_map = true
								}

								if chan_in_map {
									chan_feature := feature
									chan_feature.F_line_num = fileset.Position(x.Pos()).Line
									chan_feature.F_type = CHAN_MAP
									counter.Chan_map_count++
									counter.Features = append(counter.Features, &chan_feature)
								}
							}

						case *ast.UnaryExpr:
							switch expr := call_expr.X.(type) {
							case *ast.CompositeLit:
								switch array_type := expr.Type.(type) {
								case *ast.Ident:
									// Possible assignment of a struct struct = Struct{bla:0, bla1}
									for _, elt := range expr.Elts {
										switch valueExp := elt.(type) {
										case *ast.KeyValueExpr:
											switch ident := valueExp.Key.(type) {
											case *ast.Ident:
												switch call := valueExp.Value.(type) {
												case *ast.CallExpr:
													ident1 := value_spec.Names[index]
													checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
												}
											}
										}
									}

								case *ast.ArrayType:
									checkArrayType(array_type, counter, feature, fileset, 1)
								case *ast.MapType:
									chan_in_map := false
									// we have a declaration of a map
									switch array_type.Key.(type) {
									case *ast.ChanType:
										chan_in_map = true
									}
									switch array_type.Value.(type) {
									case *ast.ChanType:
										chan_in_map = true
									}

									if chan_in_map {
										chan_feature := feature
										chan_feature.F_line_num = fileset.Position(x.Pos()).Line
										chan_feature.F_type = CHAN_MAP
										counter.Chan_map_count++
										counter.Features = append(counter.Features, &chan_feature)
									}
								}
							}
						}
					}
				}
			}
		}
	// if it is starting a goroutine
	case *ast.GoStmt:
		go_feature := Feature{
			F_filename:     feature.F_filename,
			F_package_name: feature.F_package_name,
			F_line_num:     fileset.Position(x.Pos()).Line}
		go_feature.F_type = GOROUTINE
		counter.Go_count++
		counter.Features = append(counter.Features, &go_feature)
	// if it is sending on a channel
	case *ast.SendStmt:
		send_feature := Feature{
			F_filename:     feature.F_filename,
			F_package_name: feature.F_package_name,
			F_line_num:     fileset.Position(x.Pos()).Line}
		send_feature.F_type = SEND
		counter.Send_count++
		counter.Features = append(counter.Features, &send_feature)
	// if it is an unary expression
	case *ast.UnaryExpr:
		if x.Op.String() == "<-" {
			send_feature := Feature{
				F_filename:     feature.F_filename,
				F_package_name: feature.F_package_name,
				F_line_num:     fileset.Position(x.Pos()).Line}
			send_feature.F_type = RECEIVE
			counter.Rcv_count++
			counter.Features = append(counter.Features, &send_feature)
		}
	// if it is an assignment (different to declaration)
	case *ast.AssignStmt:
		// look for a make(chan X) or a make(chan X,n)
		for index, rh := range x.Rhs {
			switch call_expr := rh.(type) {
			case *ast.CallExpr:
				switch ident := x.Lhs[index].(type) {
				case *ast.Ident:
					checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, true)
				case *ast.SelectorExpr:
					if ident.X != nil && ident.Sel != nil {
						switch name := ident.X.(type) {
						case *ast.Ident:
							checkDepthChan(call_expr, feature, env, counter, ident.Sel.Name+"."+name.Name, fileset, true)
						}
					}
				}

			case *ast.CompositeLit:
				switch array_type := call_expr.Type.(type) {
				case *ast.Ident:
					// Possible assignment of a struct struct = Struct{bla:0, bla1}
					for _, elt := range call_expr.Elts {
						switch valueExp := elt.(type) {
						case *ast.KeyValueExpr:
							switch ident := valueExp.Key.(type) {
							case *ast.Ident:
								switch call := valueExp.Value.(type) {
								case *ast.CallExpr:
									switch ident1 := x.Lhs[index].(type) {
									case *ast.Ident:
										checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
									}
								}
							}
						}
					}

				case *ast.ArrayType:
					checkArrayType(array_type, counter, feature, fileset, 1)
				case *ast.MapType:
					chan_in_map := false
					// we have a declaration of a map
					switch array_type.Key.(type) {
					case *ast.ChanType:
						chan_in_map = true
					}
					switch array_type.Value.(type) {
					case *ast.ChanType:
						chan_in_map = true
					}

					if chan_in_map {
						chan_feature := feature
						chan_feature.F_line_num = fileset.Position(x.Pos()).Line
						chan_feature.F_type = CHAN_MAP
						counter.Chan_map_count++
						counter.Features = append(counter.Features, &chan_feature)
					}
				}

			case *ast.UnaryExpr:
				switch call_expr := call_expr.X.(type) {
				case *ast.CompositeLit:
					switch array_type := call_expr.Type.(type) {
					case *ast.Ident:
						// Possible assignment of a struct struct = Struct{bla:0, bla1}
						for _, elt := range call_expr.Elts {
							switch valueExp := elt.(type) {
							case *ast.KeyValueExpr:
								switch ident := valueExp.Key.(type) {
								case *ast.Ident:
									switch call := valueExp.Value.(type) {
									case *ast.CallExpr:
										switch ident1 := x.Lhs[index].(type) {
										case *ast.Ident:
											checkDepthChan(call, feature, env, counter, ident1.Name+"."+ident.Name, fileset, true)
										}
									}
								}
							}
						}

					case *ast.ArrayType:
						checkArrayType(array_type, counter, feature, fileset, 1)
					case *ast.MapType:
						chan_in_map := false
						// we have a declaration of a map
						switch array_type.Key.(type) {
						case *ast.ChanType:
							chan_in_map = true
						}
						switch array_type.Value.(type) {
						case *ast.ChanType:
							chan_in_map = true
						}

						if chan_in_map {
							chan_feature := feature
							chan_feature.F_line_num = fileset.Position(x.Pos()).Line
							chan_feature.F_type = CHAN_MAP
							counter.Chan_map_count++
							counter.Features = append(counter.Features, &chan_feature)
						}
					}
				}
			}
		}
	// if it is a declaration
	case *ast.DeclStmt:
		// look for a make(chan X) or a make(chan X,n)  waitgroup (var wg *sync.Waitgroup) and mutexes (var mu *sync.Mutex)
		switch decl := x.Decl.(type) {
		case *ast.GenDecl:
			// Look for declaration of a waitgroup
			if decl.Tok == token.VAR {
				for _, spec := range decl.Specs {
					switch value := spec.(type) {
					case *ast.ValueSpec:
						switch value_type := value.Type.(type) {
						case *ast.Ident:
							// 	looking for a declaration of a struct
							for index, exp := range value.Values {
								switch composite := exp.(type) {
								case *ast.CompositeLit:
									for _, elt := range composite.Elts {
										switch valueExp := elt.(type) {
										case *ast.KeyValueExpr:
											switch ident := valueExp.Key.(type) {
											case *ast.Ident:
												switch call := valueExp.Value.(type) {
												case *ast.CallExpr:
													checkDepthChan(call, feature, env, counter, value.Names[index].Name+"."+ident.Name, fileset, true)
												}
											}
										}
									}
								}
							}
						case *ast.ArrayType:
							// we have a declaration of an array
							num_of_arrays := len(value.Names)
							checkArrayType(value_type, counter, feature, fileset, num_of_arrays)

						case *ast.MapType:
							chan_in_map := false
							// we have a declaration of a map
							switch value_type.Key.(type) {
							case *ast.ChanType:
								chan_in_map = true
							}
							switch value_type.Value.(type) {
							case *ast.ChanType:
								chan_in_map = true
							}

							if chan_in_map {
								chan_feature := feature
								chan_feature.F_type = CHAN_MAP
								counter.Chan_map_count++
								counter.Features = append(counter.Features, &chan_feature)
							}
						}

						for index, val := range value.Values {
							switch call_expr := val.(type) {
							case *ast.CallExpr:
								checkDepthChan(call_expr, feature, env, counter, value.Names[index].Name, fileset, true)
							}
						}
					}
				}
			}
		}

		// Look if the type of LHS is a mutex or a waitgroup
	// if it is a for loop
	case *ast.ForStmt:
		makeChanInFor(x, feature, env, counter, fileset)
		// look in the block and see if goroutine are created in a for loop
		for _, stmt := range x.Body.List {
			switch x_node := stmt.(type) {
			case *ast.GoStmt:
				go_feature := feature
				go_feature.F_type = GO_IN_FOR
				counter.Go_in_for_count++
				go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
				counter.Features = append(counter.Features, &go_feature)
				switch bin_expr := x.Cond.(type) {
				case *ast.BinaryExpr:
					if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
						// check if the right hand side is a constant
						val, isCons := isConstant(bin_expr.Y)
						if isCons {
							go_feature := feature
							go_feature.F_type = GO_IN_CONSTANT_FOR
							go_feature.F_number = strconv.Itoa(val)
							go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
							counter.Go_in_constant_for_count++
							counter.Features = append(counter.Features, &go_feature)
						}
					} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
						// check if the initialisation is a constant
						switch assign := x.Init.(type) {
						case *ast.AssignStmt:
							for _, rh := range assign.Rhs {
								val, isCons := isConstant(rh)
								if isCons {
									go_feature := feature
									go_feature.F_type = GO_IN_CONSTANT_FOR
									go_feature.F_line_num = fileset.Position(x_node.Pos()).Line
									go_feature.F_number = strconv.Itoa(val)
									counter.Go_in_constant_for_count++
									counter.Features = append(counter.Features, &go_feature)
								}
							}
						}
					}
				}
			}
		}
	// if it is a for range loop
	case *ast.RangeStmt:
		// check if the stmt is a range over a channel

		if x.Key != nil {
			switch ident1 := x.Key.(type) {

			case *ast.Ident:
				if ident1.Obj != nil {
					switch assign := ident1.Obj.Decl.(type) {
					case *ast.AssignStmt:

						for _, rh := range assign.Rhs {
							switch unary := rh.(type) {
							case *ast.UnaryExpr:
								if unary.Op == token.RANGE {

									switch chan_type := unary.X.(type) {
									case *ast.Ident:
										// trying to range over a channel
										if chan_type.Obj != nil {
											found, _ := isChan(unary.X, env)
											if found {
												range_feature := feature
												range_feature.F_type = RANGE_OVER_CHAN
												range_feature.F_line_num = fileset.Position(unary.Pos()).Line
												counter.Range_over_chan_count++
												counter.Features = append(counter.Features, &range_feature)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			switch ident1 := x.X.(type) {
			case *ast.Ident:
				if ident1.Obj != nil {
					found, _ := isChan(ident1, env)
					if found {
						range_feature := feature
						range_feature.F_type = RANGE_OVER_CHAN
						range_feature.F_line_num = fileset.Position(ident1.Pos()).Line
						counter.Range_over_chan_count++
						counter.Features = append(counter.Features, &range_feature)
					}
				}
			}
		}
		if x.Body != nil {
			for _, stmt := range x.Body.List {
				switch assign_stmt := stmt.(type) {
				case *ast.GoStmt:
					go_in_for := Feature{
						F_filename:     filename,
						F_package_name: package_name,
						F_type:         GO_IN_FOR}
					counter.Go_in_for_count++
					go_in_for.F_line_num = fileset.Position(assign_stmt.Pos()).Line
					counter.Features = append(counter.Features, &go_in_for)

				case *ast.AssignStmt:
					for _, expr := range assign_stmt.Rhs {
						found, chan_name := isChan(expr, env)
						if found {
							assign_chan_in_for := Feature{
								F_type:         ASSIGN_CHAN_IN_FOR,
								F_filename:     filename,
								F_package_name: package_name}
							counter.Assign_chan_in_for_count++
							assign_chan_in_for.F_number = chan_name
							assign_chan_in_for.F_line_num = fileset.Position(expr.Pos()).Line
							counter.Features = append(counter.Features, &assign_chan_in_for)
						}
					}
				}
			}
		}
		makeChanInRange(x, feature, env, counter, fileset)
	// a standalone statement, like a function all, and any return type is not used
	case *ast.ExprStmt:
		// looking for a close
		switch call_expr := x.X.(type) {
		case *ast.CallExpr:
			switch ident := call_expr.Fun.(type) {
			case *ast.Ident:
				if ident.Name == "close" && len(call_expr.Args) == 1 {
					// we have a close
					found, _ := isChan(call_expr.Args[0], env)
					if found {
						// we have a close on a chan
						close_feature := feature
						counter.Close_chan_count++
						close_feature.F_type = CLOSE_CHAN
						close_feature.F_line_num = fileset.Position(ident.Pos()).Line
						counter.Features = append(counter.Features, &close_feature)
					}
				}
			}
		}
	// a select statement
	case *ast.SelectStmt:
		// that is not empty
		if x.Body != nil {
			// trackers for each case in this select
			// if another select reached, !!_!_!_! must be tracked differently
			// each kept/discard receive is tracked separately
			// timeouts are subset of either
			var with_default bool = false
			// var with_sync_send_action bool = false
			// var with_async_send_action bool = false
			// var with_sync_recv_kept_action bool = false
			// var with_async_recv_kept_action bool = false
			// var with_sync_recv_discard_action bool = false
			// var with_async_recv_discard_action bool = false
			// var with_timeout bool = false
			// // instances in this select
			// var send_sync_action_count int = 0
			// var send_async_action_count int = 0
			// var recv_sync_kept_action_count int = 0
			// var recv_async_kept_action_count int = 0
			// var recv_sync_discard_action_count int = 0
			// var recv_async_discard_action_count int = 0
			var timeout_count int = 0
			for _, stmt := range x.Body.List {
				switch comm := stmt.(type) {
				// for each case in select
				case *ast.CommClause:
					if comm.Comm == nil {
						// we have a select with a default
						with_default = true
					} else {
						switch stmt_type := comm.Comm.(type) {
						case *ast.AssignStmt:
							// receive, data kept: possible timeout
							// with_recv_kept_action = true
							// recv_kept_action_count++
							// look for source of receive :chan or timeout
							for _, rhs_stmt := range stmt_type.Rhs {
								switch unary_expr := rhs_stmt.(type) {
								case *ast.UnaryExpr:
									// both use this
									switch call_or_ident := unary_expr.X.(type) {
									case *ast.Ident:
										// a normal receive, look for chan

									case *ast.CallExpr:
										// this is a return from a function
										_with_timeout := isTimeout(call_or_ident)
										if _with_timeout {
											// with_timeout = true
											timeout_count++
										}
									}
								}
							}
							// see where value was saved
							for _, lhs_stmt := range stmt_type.Lhs {
								switch unary_expr := lhs_stmt.(type) {
								case *ast.Ident:
									lhs_obj := unary_expr.Obj
									if lhs_obj.Kind == ast.Var {
										// this is a variable assignment
									}
								}
							}
						case *ast.ExprStmt:
							// receive, data discarded: possible timeout
							// with_recv_discard_action = true
							// recv_discard_action_count++
							// look for timeout
							switch unary_expr := stmt_type.X.(type) {
							case *ast.UnaryExpr:
								// both use this
								switch call_or_ident := unary_expr.X.(type) {
								case *ast.Ident:
									// this is a normal receive
								case *ast.CallExpr:
									// this is a return from a function
									_with_timeout := isTimeout(call_or_ident)
									if _with_timeout {
										// with_timeout = true
										timeout_count++
									}
								}
							}
						case *ast.SendStmt:
							// send
							// with_send_action = true
							// send_action_count++
						}
					}
				}
			}
			select_feature := feature

			if with_default {
				// empty with just default
				select_feature.F_type = DEFAULT_SELECT
				counter.Default_select_count++
			} else {
				// empty select
				select_feature.F_type = SELECT
				counter.Select_count++
			}
			select_feature.F_number = strconv.Itoa(len(x.Body.List))
			select_feature.F_line_num = fileset.Position(x.Pos()).Line
			counter.Features = append(counter.Features, &select_feature)
		}
	case *ast.DeferStmt:
		if x.Call != nil {
			call_expr := x.Call
			switch ident := call_expr.Fun.(type) {
			case *ast.Ident:
				if ident.Name == "close" && len(call_expr.Args) == 1 {
					found, _ := isChan(call_expr.Args[0], env)
					if found {
						// we have a close on a chan
						close_feature := feature
						counter.Close_chan_count++
						close_feature.F_type = CLOSE_CHAN
						close_feature.F_line_num = fileset.Position(call_expr.Pos()).Line
						counter.Features = append(counter.Features, &close_feature)
					}
				}
			}
		}
	case *ast.FuncDecl:
		// look for a <-chan, chan<- or chan as function fields
		for _, field := range x.Type.Params.List {
			switch chan_type := field.Type.(type) {
			case *ast.ChanType:
				switch chan_type.Dir {
				case ast.RECV:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = RECEIVE_CHAN
					counter.Receive_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				case ast.SEND:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = SEND_CHAN
					counter.Send_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				default:
					chan_feature := feature
					chan_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					chan_feature.F_type = PARAM_CHAN
					counter.Param_chan_count++
					counter.Features = append(counter.Features, &chan_feature)
				}
			}
		}
	}
}

func makeChanInFor(forStmt *ast.ForStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet) {
	for _, block := range forStmt.Body.List {
		switch stmt := block.(type) {
		case *ast.AssignStmt:
			// chan in for
			for index, rh := range stmt.Rhs {
				switch call_expr := rh.(type) {
				case *ast.CallExpr:
					switch ident := stmt.Lhs[index].(type) {
					case *ast.Ident:
						if checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, false) {
							switch bin_expr := forStmt.Cond.(type) {
							case *ast.BinaryExpr:
								if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
									// check if the right hand side is a constant
									val, isCons := isConstant(bin_expr.Y)
									if isCons {
										make_chan_in_for := Feature{}
										make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
										if stmt.Tok == token.DEFINE {
											make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
										}
										make_chan_in_for.F_filename = feature.F_filename
										make_chan_in_for.F_package_name = feature.F_package_name
										make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
										make_chan_in_for.F_number = strconv.Itoa(val)
										counter.Make_chan_in_constant_for_count++
										counter.Features = append(counter.Features, &make_chan_in_for)
									}
									// }
								} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
									// check if the initialisation is a constant
									switch assign := forStmt.Init.(type) {
									case *ast.AssignStmt:
										for _, rh := range assign.Rhs {
											val, isCons := isConstant(rh)
											if isCons {
												make_chan_in_for := Feature{}
												make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
												if stmt.Tok == token.DEFINE {
													make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
												}
												make_chan_in_for.F_filename = feature.F_filename
												make_chan_in_for.F_package_name = feature.F_package_name
												make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
												make_chan_in_for.F_number = strconv.Itoa(val)
												counter.Make_chan_in_constant_for_count++
												counter.Features = append(counter.Features, &make_chan_in_for)
											}
										}
									}
								} else {
									make_chan_in_for := Feature{}
									make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
									if stmt.Tok == token.DEFINE {
										make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
									}
									make_chan_in_for.F_filename = feature.F_filename
									make_chan_in_for.F_package_name = feature.F_package_name
									make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
									counter.Make_chan_in_constant_for_count++
									counter.Features = append(counter.Features, &make_chan_in_for)
								}
							default:
								make_chan_in_for := Feature{}
								make_chan_in_for.F_type = ASSIGN_CHAN_IN_FOR
								if stmt.Tok == token.DEFINE {
									make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
								}
								make_chan_in_for.F_filename = feature.F_filename
								make_chan_in_for.F_package_name = feature.F_package_name
								make_chan_in_for.F_line_num = fileset.Position(ident.Pos()).Line
								counter.Make_chan_in_constant_for_count++
								counter.Features = append(counter.Features, &make_chan_in_for)
							}
						}
					}

				}
			}

		case *ast.DeclStmt: // is the declaration in a constant or not for loop ?
			if chanDecleration(stmt, feature, env, counter, fileset, false) {
				switch bin_expr := forStmt.Cond.(type) {
				case *ast.BinaryExpr:
					if bin_expr.Op == token.LEQ || bin_expr.Op == token.LSS { // <, <=
						// check if the right hand side is a constant
						val, isCons := isConstant(bin_expr.Y)
						if isCons {
							make_chan_in_for := Feature{}
							make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
							make_chan_in_for.F_filename = feature.F_filename
							make_chan_in_for.F_package_name = feature.F_package_name
							make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
							make_chan_in_for.F_number = strconv.Itoa(val)
							counter.Make_chan_in_constant_for_count++
							counter.Features = append(counter.Features, &make_chan_in_for)
						}
					} else if bin_expr.Op == token.GEQ || bin_expr.Op == token.GTR { // >, >=
						// check if the initialisation is a constant
						switch assign := forStmt.Init.(type) {
						case *ast.AssignStmt:
							for _, rh := range assign.Rhs {
								val, isCons := isConstant(rh)
								if isCons {
									make_chan_in_for := Feature{}
									make_chan_in_for.F_type = MAKE_CHAN_IN_CONSTANT_FOR
									make_chan_in_for.F_filename = feature.F_filename
									make_chan_in_for.F_package_name = feature.F_package_name
									make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
									make_chan_in_for.F_number = strconv.Itoa(val)
									counter.Make_chan_in_constant_for_count++
									counter.Features = append(counter.Features, &make_chan_in_for)
								}
							}
						}
					} else {
						make_chan_in_for := Feature{}
						make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
						make_chan_in_for.F_filename = feature.F_filename
						make_chan_in_for.F_package_name = feature.F_package_name
						make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
						counter.Make_chan_in_constant_for_count++
						counter.Features = append(counter.Features, &make_chan_in_for)
					}
				default:
					make_chan_in_for := Feature{}
					make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
					make_chan_in_for.F_filename = feature.F_filename
					make_chan_in_for.F_package_name = feature.F_package_name
					make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
					counter.Make_chan_in_constant_for_count++
					counter.Features = append(counter.Features, &make_chan_in_for)
				}
			}
		}
	}

	for _, stmt := range forStmt.Body.List {
		switch x_node := stmt.(type) {

		case *ast.AssignStmt:
			for _, expr := range x_node.Rhs {
				found, chan_name := isChan(expr, env)
				chan_feature := feature
				if found {
					chan_feature.F_type = ASSIGN_CHAN_IN_FOR
					if x_node.Tok == token.DEFINE {
						chan_feature.F_type = MAKE_CHAN_IN_CONSTANT_FOR
					}
					chan_feature.F_line_num = fileset.Position(expr.Pos()).Line
					chan_feature.F_number = chan_name
					counter.Assign_chan_in_for_count++
					counter.Features = append(counter.Features, &chan_feature)
				}
			}
		}
	}
}

func makeChanInRange(rangeStmt *ast.RangeStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet) {

	for _, block := range rangeStmt.Body.List {
		switch stmt := block.(type) {
		case *ast.AssignStmt:
			// chan in for
			for index, rh := range stmt.Rhs {
				switch call_expr := rh.(type) {
				case *ast.CallExpr:
					switch ident := stmt.Lhs[index].(type) {
					case *ast.Ident:
						if checkDepthChan(call_expr, feature, env, counter, ident.Name, fileset, false) {
							make_chan_in_for := Feature{}
							make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
							make_chan_in_for.F_filename = feature.F_filename
							make_chan_in_for.F_package_name = feature.F_package_name
							make_chan_in_for.F_line_num = fileset.Position(call_expr.Pos()).Line
							counter.Make_chan_in_constant_for_count++
							counter.Features = append(counter.Features, &make_chan_in_for)
						}
					}
				}
			}

		case *ast.DeclStmt:

			if chanDecleration(stmt, feature, env, counter, fileset, false) {
				make_chan_in_for := Feature{}
				make_chan_in_for.F_type = MAKE_CHAN_IN_FOR
				make_chan_in_for.F_filename = feature.F_filename
				make_chan_in_for.F_package_name = feature.F_package_name
				make_chan_in_for.F_line_num = fileset.Position(stmt.Pos()).Line
				counter.Make_chan_in_constant_for_count++
				counter.Features = append(counter.Features, &make_chan_in_for)
			}
		}
	}
}

func chanDecleration(stmt *ast.DeclStmt, feature Feature, env *[]string, counter *Counter, fileset *token.FileSet, add bool) bool {
	var found_decl bool = false
	switch decl := stmt.Decl.(type) {
	case *ast.GenDecl:
		if decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				switch value := spec.(type) {
				case *ast.ValueSpec:
					switch value.Type.(type) {
					case *ast.ChanType:
						// we have a var x chan X
						if value.Values != nil {
							if len(value.Values) == len(value.Names) {
								for index, val := range value.Values {
									switch call_expr := val.(type) {
									case *ast.CallExpr:
										found_decl = checkDepthChan(call_expr, feature, env, counter, value.Names[index].Name, fileset, add)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return found_decl
}

func checkArrayType(array_type *ast.ArrayType, counter *Counter, feature Feature, fileset *token.FileSet, num_of_arrays int) {
	switch chan_type := array_type.Elt.(type) {
	case *ast.ChanType:
		//we have an array of chan
		if array_type.Len != nil {
			// check if constant
			val, isCons := isConstant(array_type.Len)
			if isCons {
				for i := 0; i < num_of_arrays; i++ {
					array_feature := feature
					array_feature.F_type = CONSTANT_CHAN_ARRAY
					array_feature.F_number = strconv.Itoa(val)
					array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					counter.Constant_chan_array_count += num_of_arrays
					counter.Features = append(counter.Features, &array_feature)
				}
			} else {
				for i := 0; i < num_of_arrays; i++ {
					array_feature := feature
					array_feature.F_type = ARRAY_OF_CHANNELS
					counter.Array_of_channels_count += num_of_arrays
					array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
					counter.Features = append(counter.Features, &array_feature)
				}
			}
		} else {
			for i := 0; i < num_of_arrays; i++ {
				array_feature := feature
				array_feature.F_type = CHAN_SLICE
				counter.Chan_slice_count += num_of_arrays
				array_feature.F_line_num = fileset.Position(chan_type.Pos()).Line
				counter.Features = append(counter.Features, &array_feature)
			}
		}
	}
}

func checkDepthChan(call_expr *ast.CallExpr, feature Feature, env *[]string, counter *Counter, chan_name string, fileset *token.FileSet, add bool) bool {
	var chan_found bool = false
	switch ident := call_expr.Fun.(type) {
	case *ast.Ident:
		if ident.Name == "make" {
			if len(call_expr.Args) > 0 {
				switch chan_type := call_expr.Args[0].(type) {
				case *ast.ChanType:
					chan_found = true
					*env = append(*env, chan_name)
					switch chan_type.Value.(type) {

					case *ast.ChanType:
						chan_feature := Feature{
							F_filename:     feature.F_filename,
							F_package_name: feature.F_package_name,
							F_line_num:     fileset.Position(call_expr.Pos()).Line}
						chan_feature.F_type = CHAN_OF_CHANS
						chan_feature.F_number = chan_name
						counter.Chan_of_chans_count++
						counter.Features = append(counter.Features, &chan_feature)
					default:
						if len(call_expr.Args) > 1 {
							val, isCons := isConstant(call_expr.Args[1])
							if isCons {
								if add {
									if val != 0 {
										chan_feature := Feature{
											F_filename:     feature.F_filename,
											F_package_name: feature.F_package_name,
											F_line_num:     fileset.Position(call_expr.Pos()).Line}
										chan_feature.F_type = KNOWN_CHAN_DEPTH
										chan_feature.F_number = strconv.Itoa(val)
										counter.Known_chan_depth_count++
										counter.Features = append(counter.Features, &chan_feature)
									} else {
										chan_feature := Feature{
											F_filename:     feature.F_filename,
											F_package_name: feature.F_package_name,
											F_line_num:     fileset.Position(call_expr.Pos()).Line}
										chan_feature.F_type = MAKE_CHAN
										counter.Sync_Chan_count++
										counter.Features = append(counter.Features, &chan_feature)
									}
								}
							} else {
								if add {
									chan_feature := Feature{
										F_filename:     feature.F_filename,
										F_package_name: feature.F_package_name,
										F_line_num:     fileset.Position(call_expr.Pos()).Line}
									chan_feature.F_type = UNKNOWN_CHAN_DEPTH //unknown depth
									counter.Unknown_chan_depth_count++
									counter.Features = append(counter.Features, &chan_feature)
								}
							}
						} else {
							if add {
								chan_feature := Feature{
									F_filename:     feature.F_filename,
									F_package_name: feature.F_package_name,
									F_line_num:     fileset.Position(call_expr.Pos()).Line}
								chan_feature.F_type = MAKE_CHAN
								counter.Sync_Chan_count++
								counter.Features = append(counter.Features, &chan_feature)
							}
						}
					}
				}
			}
		}
	}

	return chan_found
}

func isTimeout(callExpr *ast.CallExpr) bool {
	switch sel_expr := callExpr.Fun.(type) {
	// found function
	case *ast.SelectorExpr:
		var sel_expr_x_name bool
		var sel_expr_sel_name bool
		switch sel_expr_x := sel_expr.X.(type) {
		case *ast.Ident:
			if sel_expr_x.Name == "time" {
				sel_expr_x_name = true
			}
		}
		if sel_expr.Sel.Name == "After" {
			sel_expr_sel_name = true
		}
		// check if timeout found
		if sel_expr_x_name && sel_expr_sel_name {
			return true
		}
	}
	return false
}

func isConstant(node ast.Node) (int, bool) {
	var isCons bool = false
	var value int = 0
	switch ident := node.(type) {
	case *ast.Ident:
		if ident.Obj != nil {
			if ident.Obj.Kind == ast.Con {
				switch value_spec := ident.Obj.Decl.(type) {
				case *ast.ValueSpec:

					if value_spec.Values != nil && len(value_spec.Values) > 0 {
						switch val := value_spec.Values[0].(type) {
						case *ast.BasicLit:
							parsed_int, _ := strconv.Atoi(val.Value)
							value = int(parsed_int)
							isCons = true
						case *ast.Ident:
							value, isCons = isConstant(val)
						}
					}
				}
			}
		}
	case *ast.BasicLit:
		if ident.Kind == token.INT {
			isCons = true
			parsed_int, _ := strconv.Atoi(ident.Value)
			value = int(parsed_int)
		}
	default:
		isCons = false
	}

	return value, isCons
}

func chanType() {

}

func isChan(node interface{}, env *[]string) (bool, string) {

	chan_name := ""
	switch make_chan := node.(type) {
	case *ast.AssignStmt:
		var chan_found bool = false
		ast.Inspect(make_chan, func(x_node ast.Node) bool {
			switch x_node.(type) {
			case *ast.ChanType:
				chan_found = true

				return false
			}
			return true
		})

		if !chan_found {
			for _, rh := range make_chan.Rhs {
				switch ident := rh.(type) {
				case *ast.Ident:
					for _, name := range *env {
						if name == ident.Name {
							chan_found = true
							chan_name = name
							break
						}
					}
				}
			}
		}
		return chan_found, chan_name
	case *ast.Ident:
		for _, name := range *env {
			if name == make_chan.Name {
				chan_name = name
				return true, chan_name
			}
		}
	}

	return false, chan_name
}

// adds all channels globally accessible in a file to env
func addGlobalVarToEnv(file *ast.File, env *[]string) {
	for _, decl := range file.Decls {
		// for every global declaration in this file ( )
		switch genDecl := decl.(type) {
		case *ast.GenDecl:
			// if the declaration type is generic
			if genDecl.Tok == token.VAR {
				// if it is a variable declaration
				for _, spec := range genDecl.Specs {
					switch value_spec := spec.(type) {
					case *ast.ValueSpec:
						// if it is a constant or a variable
						for index, value := range value_spec.Values {
							switch call_expr := value.(type) {
							case *ast.CallExpr:
								// if it has arguments, make(chan, ...)
								switch ident := call_expr.Fun.(type) {
								// get its identifier from its function expression
								case *ast.Ident:
									// if it is creating a channel
									if ident.Name == "make" {
										// if it is a valid declaration
										if len(call_expr.Args) > 0 {
											// check if it is sync or async
											if len(call_expr.Args) > 1 {

											}
											switch call_expr.Args[0].(type) {
											// if it is a channel
											case *ast.ChanType:
												*env = append(*env, value_spec.Names[index].Name)
											}
										}
									}
								}
								// case *ast.Ident:
								// 	// var or const
								// 	var_type := call_expr.Obj.Kind
								// 	data_type := call_expr.k

							}
						}
					}
				}
			}
		}
	}
}
