package main

type funcTweak struct {
	// name specifies the name of the Go function to be tweaked.
	name string

	// copy copies all the definitions for this function tweak from the named
	// function. Templates are parsed under the new context.
	copy string

	// params specifies a map of zero or more tweaks for specific parameters.
	params paramTweaks

	// result defines the function result as presented at the end of the func line.
	// Simple type changes are handled automatically. More involved multi-value
	// results will require an appropriate after snippet to handle the return.
	result string

	// before is a code snippet to be injected before the C function call.
	// It may use the following template variables and functions:
	//
	//                          . - dot holds the Func being tweaked
	//         {{copyDoc "Func"}} - replaced by the respective function documentation
	//  {{paramGoType . "param"}} - replaced by the respective parameter Go type
	//
	before string

	// after is a code snippet to be injected after the C function call.
	// It may use the same template functions as available for before.
	after string

	// doc defines the documentation for the function. It may use the same
	// template functions as available for before and after.
	doc string
}

type paramTweak struct {
	// rename changes the parameter name in the Go function while keeping the C
	// function call unchanged. The before snippet must define a proper variable
	// to be used under the original name.
	rename string

	// replace changes the parameter name in the C function call to a variable
	// named "<original name>_c", while keeping the Go parameter name unchanged.
	// The before and after snippets must manipulate the two values as needed.
	replace bool

	// retype changes the Go parameter type.
	retype string

	// output flags the parameter as an output parameter, which causes it to be
	// omitted from the input parameter list and added to the result list.
	output bool

	// unnamed causes the name of a result parameter to be omitted if possible.
	unnamed bool

	// single flags the parameter as carrying a single value rather than a slice,
	// when the parameter is originally defined as a pointer.
	single bool

	// omit drops the parameter from the Go function. The before snippet must
	// define a variable with the proper name for the C function call to use.
	omit bool
}

type paramTweaks map[string]paramTweak

var paramNameFixes = map[string]string{
	"binaryformat":   "binaryFormat",
	"bufsize":        "bufSize",
	"indx":           "index",
	"infolog":        "infoLog",
	"internalformat": "internalFormat",
	"precisiontype":  "precisionType",
}

var funcTweakList = []funcTweak{{
	name: "Accum",
	doc: `
		executes an operation on the accumulation buffer.

		Parameter op defines the accumulation buffer operation (GL.ACCUM, GL.LOAD,
		GL.ADD, GL.MULT, or GL.RETURN) and specifies how the value parameter is
		used.

		The accumulation buffer is an extended-range color buffer. Images are not
		rendered into it. Rather, images rendered into one of the color buffers
		are added to the contents of the accumulation buffer after rendering.
		Effects such as antialiasing (of points, lines, and polygons), motion
		blur, and depth of field can be created by accumulating images generated
		with different transformation matrices.

		Each pixel in the accumulation buffer consists of red, green, blue, and
		alpha values. The number of bits per component in the accumulation buffer
		depends on the implementation. You can examine this number by calling
		GetIntegerv four times, with arguments GL.ACCUM_RED_BITS,
		GL.ACCUM_GREEN_BITS, GL.ACCUM_BLUE_BITS, and GL.ACCUM_ALPHA_BITS.
		Regardless of the number of bits per component, the range of values stored
		by each component is (-1, 1). The accumulation buffer pixels are mapped
		one-to-one with frame buffer pixels.

		All accumulation buffer operations are limited to the area of the current
		scissor box and applied identically to the red, green, blue, and alpha
		components of each pixel. If a Accum operation results in a value outside
		the range (-1, 1), the contents of an accumulation buffer pixel component
		are undefined.

		The operations are as follows:

		  GL.ACCUM
		      Obtains R, G, B, and A values from the buffer currently selected for
		      reading (see ReadBuffer). Each component value is divided by 2 n -
		      1 , where n is the number of bits allocated to each color component
		      in the currently selected buffer. The result is a floating-point
		      value in the range 0 1 , which is multiplied by value and added to
		      the corresponding pixel component in the accumulation buffer,
		      thereby updating the accumulation buffer.

		  GL.LOAD
		      Similar to GL.ACCUM, except that the current value in the
		      accumulation buffer is not used in the calculation of the new value.
		      That is, the R, G, B, and A values from the currently selected
		      buffer are divided by 2 n - 1 , multiplied by value, and then stored
		      in the corresponding accumulation buffer cell, overwriting the
		      current value.

		  GL.ADD
		      Adds value to each R, G, B, and A in the accumulation buffer.

		  GL.MULT
		      Multiplies each R, G, B, and A in the accumulation buffer by value
		      and returns the scaled component to its corresponding accumulation
		      buffer location.

		  GL.RETURN
		      Transfers accumulation buffer values to the color buffer or buffers
		      currently selected for writing. Each R, G, B, and A component is
		      multiplied by value, then multiplied by 2 n - 1 , clamped to the
		      range 0 2 n - 1 , and stored in the corresponding display buffer
		      cell. The only fragment operations that are applied to this transfer
		      are pixel ownership, scissor, dithering, and color writemasks.

		To clear the accumulation buffer, call ClearAccum with R, G, B, and A
		values to set it to, then call Clear with the accumulation buffer
		enabled.

		Error GL.INVALID_ENUM is generated if op is not an accepted value.  
		GL.INVALID_OPERATION is generated if there is no accumulation buffer.
		GL.INVALID_OPERATION is generated if Accum is executed between the
		execution of Begin and the corresponding execution of End.
	`,
}, {
	name: "AttachShader",
	doc: `
		attaches a shader object to a program object.

		In order to create an executable, there must be a way to specify the list
		of things that will be linked together. Program objects provide this
		mechanism. Shaders that are to be linked together in a program object must
		first be attached to that program object. This indicates that shader will
		be included in link operations that will be performed on program.

		All operations that can be performed on a shader object are valid whether
		or not the shader object is attached to a program object. It is
		permissible to attach a shader object to a program object before source
		code has been loaded into the shader object or before the shader object
		has been compiled. It is permissible to attach multiple shader objects of
		the same type because each may contain a portion of the complete shader.
		It is also permissible to attach a shader object to more than one program
		object. If a shader object is deleted while it is attached to a program
		object, it will be flagged for deletion, and deletion will not occur until
		DetachShader is called to detach it from all program objects to which it
		is attached.

		Error GL.INVALID_VALUE is generated if either program or shader is not a
		value generated by OpenGL. GL.INVALID_OPERATION is generated if program
		is not a program object. GL.INVALID_OPERATION is generated if shader is
		not a shader object. GL.INVALID_OPERATION is generated if shader is
		already attached to program. GL.INVALID_OPERATION is generated if
		AttachShader is executed between the execution of Begin and the
		corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "BindAttribLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	doc: `
		associates a user-defined attribute variable in the program
		object specified by program with a generic vertex attribute index. The name
		parameter specifies the name of the vertex shader attribute variable to
		which index is to be bound. When program is made part of the current state,
		values provided via the generic vertex attribute index will modify the
		value of the user-defined attribute variable specified by name.

		If name refers to a matrix attribute variable, index refers to the first
		column of the matrix. Other matrix columns are then automatically bound to
		locations index+1 for a matrix of type mat2; index+1 and index+2 for a
		matrix of type mat3; and index+1, index+2, and index+3 for a matrix of
		type mat4.

		This command makes it possible for vertex shaders to use descriptive names
		for attribute variables rather than generic variables that are numbered
		from 0 to GL.MAX_VERTEX_ATTRIBS-1. The values sent to each generic
		attribute index are part of current state, just like standard vertex
		attributes such as color, normal, and vertex position. If a different
		program object is made current by calling UseProgram, the generic vertex
		attributes are tracked in such a way that the same values will be observed
		by attributes in the new program object that are also bound to index.

		Attribute variable name-to-generic attribute index bindings for a program
		object can be explicitly assigned at any time by calling
		BindAttribLocation. Attribute bindings do not go into effect until
		LinkProgram is called. After a program object has been linked
		successfully, the index values for generic attributes remain fixed (and
		their values can be queried) until the next link command occurs.

		Applications are not allowed to bind any of the standard OpenGL vertex
		attributes using this command, as they are bound automatically when
		needed. Any attribute binding that occurs after the program object has
		been linked will not take effect until the next time the program object is
		linked.

		If name was bound previously, that information is lost. Thus you cannot
		bind one user-defined attribute variable to multiple indices, but you can
		bind multiple user-defined attribute variables to the same index.

		Applications are allowed to bind more than one user-defined attribute
		variable to the same generic vertex attribute index. This is called
		aliasing, and it is allowed only if just one of the aliased attributes is
		active in the executable program, or if no path through the shader
		consumes more than one attribute of a set of attributes aliased to the
		same location. The compiler and linker are allowed to assume that no
		aliasing is done and are free to employ optimizations that work only in
		the absence of aliasing. OpenGL implementations are not required to do
		error checking to detect aliasing. Because there is no way to bind
		standard attributes, it is not possible to alias generic attributes with
		conventional ones (except for generic attribute 0).

		BindAttribLocation can be called before any vertex shader objects are
		bound to the specified program object. It is also permissible to bind a
		generic attribute index to an attribute variable name that is never used
		in a vertex shader.

		Active attributes that are not explicitly bound will be bound by the
		linker when LinkProgram is called. The locations assigned can be queried
		by calling GetAttribLocation.

		Error GL.INVALID_VALUE is generated if index is greater than or equal to
		GL.MAX_VERTEX_ATTRIBS.
		GL.INVALID_OPERATION is generated if name starts with the reserved prefix "gl_".
		GL.INVALID_VALUE is generated if program is not a value generated by OpenGL.
		GL.INVALID_OPERATION is generated if program is not a program object.
		GL.INVALID_OPERATION is generated if BindAttribLocation is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "BindBuffer",
	doc: `
		creates or puts in use a named buffer object.
		Calling BindBuffer with target set to GL.ARRAY_BUFFER,
		GL.ELEMENT_ARRAY_BUFFER, GL.PIXEL_PACK_BUFFER or GL.PIXEL_UNPACK_BUFFER
		and buffer set to the name of the new buffer object binds the buffer
		object name to the target. When a buffer object is bound to a target, the
		previous binding for that target is automatically broken.
		
		Buffer object names are unsigned integers. The value zero is reserved, but
		there is no default buffer object for each buffer object target. Instead,
		buffer set to zero effectively unbinds any buffer object previously bound,
		and restores client memory usage for that buffer object target. Buffer
		object names and the corresponding buffer object contents are local to the
		shared display-list space (see XCreateContext) of the current GL rendering
		context; two rendering contexts share buffer object names only if they
		also share display lists.
		
		GenBuffers may be called to generate a set of new buffer object names.
		
		The state of a buffer object immediately after it is first bound is an
		unmapped zero-sized memory buffer with GL.READ_WRITE access and
		GL.STATIC_DRAW usage.
		
		While a non-zero buffer object name is bound, GL operations on the target
		to which it is bound affect the bound buffer object, and queries of the
		target to which it is bound return state from the bound buffer object.
		While buffer object name zero is bound, as in the initial state, attempts
		to modify or query state on the target to which it is bound generates an
		GL.INVALID_OPERATION error.
		
		When vertex array pointer state is changed, for example by a call to
		NormalPointer, the current buffer object binding (GL.ARRAY_BUFFER_BINDING)
		is copied into the corresponding client state for the vertex array type
		being changed, for example GL.NORMAL_ARRAY_BUFFER_BINDING. While a
		non-zero buffer object is bound to the GL.ARRAY_BUFFER target, the vertex
		array pointer parameter that is traditionally interpreted as a pointer to
		client-side memory is instead interpreted as an offset within the buffer
		object measured in basic machine units.
		
		While a non-zero buffer object is bound to the GL.ELEMENT_ARRAY_BUFFER
		target, the indices parameter of DrawElements, DrawRangeElements, or
		MultiDrawElements that is traditionally interpreted as a pointer to
		client-side memory is instead interpreted as an offset within the buffer
		object measured in basic machine units.
		
		While a non-zero buffer object is bound to the GL.PIXEL_PACK_BUFFER
		target, the following commands are affected: GetCompressedTexImage,
		GetConvolutionFilter, GetHistogram, GetMinmax, GetPixelMap,
		GetPolygonStipple, GetSeparableFilter, GetTexImage, and ReadPixels. The
		pointer parameter that is traditionally interpreted as a pointer to
		client-side memory where the pixels are to be packed is instead
		interpreted as an offset within the buffer object measured in basic
		machine units.
		
		While a non-zero buffer object is bound to the GL.PIXEL_UNPACK_BUFFER
		target, the following commands are affected: Bitmap, ColorSubTable,
		ColorTable, CompressedTexImage1D, CompressedTexImage2D,
		CompressedTexImage3D, CompressedTexSubImage1D, CompressedTexSubImage2D,
		CompressedTexSubImage3D, ConvolutionFilter1D, ConvolutionFilter2D,
		DrawPixels, PixelMap, PolygonStipple, SeparableFilter2D, TexImage1D,
		TexImage2D, TexImage3D, TexSubImage1D, TexSubImage2D, and TexSubImage3D.
		The pointer parameter that is traditionally interpreted as a pointer to
		client-side memory from which the pixels are to be unpacked is instead
		interpreted as an offset within the buffer object measured in basic
		machine units.
		
		A buffer object binding created with BindBuffer remains active until a
		different buffer object name is bound to the same target, or until the
		bound buffer object is deleted with DeleteBuffers.
		
		Once created, a named buffer object may be re-bound to any target as often
		as needed. However, the GL implementation may make choices about how to
		optimize the storage of a buffer object based on its initial binding
		target.
		
		Error GL.INVALID_ENUM is generated if target is not one of the allowable
		values.  GL.INVALID_OPERATION is generated if BindBuffer is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "1.5+"}}
	`,
}, {
	name: "BufferData",
	before: `
		if data != nil {
			size = int(data_v.Type().Size()) * data_v.Len()
		}
	`,
	doc: `
		creates a new data store for the buffer object currently
		bound to target. Any pre-existing data store is deleted. The new data
		store is created with the specified size in bytes and usage. If data is
		not nil, it must be a slice that is used to initialize the data store.
		In that case the size parameter is ignored and the store size will match
		the slice data size.

		In its initial state, the new data store is not mapped, it has a NULL
		mapped pointer, and its mapped access is GL.READ_WRITE.
		
		The target constant must be one of GL.ARRAY_BUFFER, GL.COPY_READ_BUFFER,
		GL.COPY_WRITE_BUFFER, GL.ELEMENT_ARRAY_BUFFER, GL.PIXEL_PACK_BUFFER,
		GL.PIXEL_UNPACK_BUFFER, GL.TEXTURE_BUFFER, GL.TRANSFORM_FEEDBACK_BUFFER,
		or GL.UNIFORM_BUFFER.

		The usage parameter is a hint to the GL implementation as to how a buffer
		object's data store will be accessed. This enables the GL implementation
		to make more intelligent decisions that may significantly impact buffer
		object performance. It does not, however, constrain the actual usage of
		the data store. usage can be broken down into two parts: first, the
		frequency of access (modification and usage), and second, the nature of
		that access.

		A usage frequency of STREAM and nature of DRAW is specified via the
		constant GL.STREAM_DRAW, for example.
		
		The usage frequency of access may be one of:
		
		  STREAM
		      The data store contents will be modified once and used at most a few times.
		
		  STATIC
		      The data store contents will be modified once and used many times.
		
		  DYNAMIC
		      The data store contents will be modified repeatedly and used many times.
		
		The usage nature of access may be one of:
		
		  DRAW
		      The data store contents are modified by the application, and used as
		      the source for GL drawing and image specification commands.
		
		  READ
		      The data store contents are modified by reading data from the GL,
		      and used to return that data when queried by the application.
		
		  COPY
		      The data store contents are modified by reading data from the GL,
		      and used as the source for GL drawing and image specification
		      commands.

		Clients must align data elements consistent with the requirements of the
		client platform, with an additional base-level requirement that an offset
		within a buffer to a datum comprising N bytes be a multiple of N.

		Error GL.INVALID_ENUM is generated if target is not one of the accepted
		buffer targets.  GL.INVALID_ENUM is generated if usage is not
		GL.STREAM_DRAW, GL.STREAM_READ, GL.STREAM_COPY, GL.STATIC_DRAW,
		GL.STATIC_READ, GL.STATIC_COPY, GL.DYNAMIC_DRAW, GL.DYNAMIC_READ, or
		GL.DYNAMIC_COPY.  GL.INVALID_VALUE is generated if size is negative.
		GL.INVALID_OPERATION is generated if the reserved buffer object name 0 is
		bound to target.  GL.OUT_OF_MEMORY is generated if the GL is unable to
		create a data store with the specified size.
	`,
}, {
	name: "DepthRange",
	doc: `
		specifies the mapping of depth values from normalized device
		coordinates to window coordinates.

		Parameter nearVal specifies the mapping of the near clipping plane to window
		coordinates (defaults to 0), while farVal specifies the mapping of the far
		clipping plane to window coordinates (defaults to 1).

		After clipping and division by w, depth coordinates range from -1 to 1,
		corresponding to the near and far clipping planes. DepthRange specifies a
		linear mapping of the normalized depth coordinates in this range to window
		depth coordinates. Regardless of the actual depth buffer implementation,
		window coordinate depth values are treated as though they range from 0 through 1
		(like color components). Thus, the values accepted by DepthRange are both
		clamped to this range before they are accepted.

		The default setting of (0, 1) maps the near plane to 0 and the far plane to 1.
		With this mapping, the depth buffer range is fully utilized.

		It is not necessary that nearVal be less than farVal. Reverse mappings such as
		nearVal 1, and farVal 0 are acceptable.

		GL.INVALID_OPERATION is generated if DepthRange is executed between the
		execution of Begin and the corresponding execution of End.
	`,
}, {
	name: "CompileShader",
	doc: `
		compiles the source code strings that have been stored in
		the shader object specified by shader.

		The compilation status will be stored as part of the shader object's
		state. This value will be set to GL.TRUE if the shader was compiled without
		errors and is ready for use, and GL.FALSE otherwise. It can be queried by
		calling GetShaderiv with arguments shader and GL.COMPILE_STATUS.

		Compilation of a shader can fail for a number of reasons as specified by
		the OpenGL Shading Language Specification. Whether or not the compilation
		was successful, information about the compilation can be obtained from the
		shader object's information log by calling GetShaderInfoLog.

		Error GL.INVALID_VALUE is generated if shader is not a value generated by
		OpenGL.  GL.INVALID_OPERATION is generated if shader is not a shader
		object.  GL.INVALID_OPERATION is generated if CompileShader is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name:   "CreateProgram",
	result: "glbase.Program",
	doc: `
		creates an empty program object and returns a non-zero
		value by which it can be referenced. A program object is an object to
		which shader objects can be attached. This provides a mechanism to specify
		the shader objects that will be linked to create a program. It also
		provides a means for checking the compatibility of the shaders that will
		be used to create a program (for instance, checking the compatibility
		between a vertex shader and a fragment shader). When no longer needed as
		part of a program object, shader objects can be detached.

		One or more executables are created in a program object by successfully
		attaching shader objects to it with AttachShader, successfully compiling
		the shader objects with CompileShader, and successfully linking the
		program object with LinkProgram. These executables are made part of
		current state when UseProgram is called. Program objects can be deleted
		by calling DeleteProgram. The memory associated with the program object
		will be deleted when it is no longer part of current rendering state for
		any context.

		Like display lists and texture objects, the name space for program objects
		may be shared across a set of contexts, as long as the server sides of the
		contexts share the same address space. If the name space is shared across
		contexts, any attached objects and the data associated with those attached
		objects are shared as well.

		Applications are responsible for providing the synchronization across API
		calls when objects are accessed from different execution threads.

		This function returns 0 if an error occurs creating the program object.

		Error GL.INVALID_OPERATION is generated if CreateProgram is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name:   "CreateShader",
	result: "glbase.Shader",
	doc: `
		creates an empty shader object and returns a non-zero value
		by which it can be referenced. A shader object is used to maintain the
		source code strings that define a shader. shaderType indicates the type of
		shader to be created.
		
		Two types of shaders are supported. A shader of type GL.VERTEX_SHADER is a
		shader that is intended to run on the programmable vertex processor and
		replace the fixed functionality vertex processing in OpenGL. A shader of
		type GL.FRAGMENT_SHADER is a shader that is intended to run on the
		programmable fragment processor and replace the fixed functionality
		fragment processing in OpenGL.
		
		When created, a shader object's GL.SHADER_TYPE parameter is set to either
		GL.VERTEX_SHADER or GL.FRAGMENT_SHADER, depending on the value of
		shaderType.

		Like display lists and texture objects, the name space for shader objects
		may be shared across a set of contexts, as long as the server sides of the
		contexts share the same address space. If the name space is shared across
		contexts, any attached objects and the data associated with those attached
		objects are shared as well.
		
		This function returns 0 if an error occurs creating the shader object.
		
		Error GL.INVALID_ENUM is generated if shaderType is not an accepted value.
		GL.INVALID_OPERATION is generated if CreateShader is executed between the
		execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GenBuffers",
	params: paramTweaks{
		"buffers": {output: true, unnamed: true, retype: "[]glbase.Buffer"},
	},
	before: `
		buffers := make([]glbase.Buffer, n)
	`,
	doc: `
		returns n buffer object names. There is no guarantee that
		the names form a contiguous set of integers; however, it is guaranteed
		that none of the returned names was in use immediately before the call to
		GenBuffers.

		Buffer object names returned by a call to GenBuffers are not returned by
		subsequent calls, unless they are first deleted with DeleteBuffers.

		No buffer objects are associated with the returned buffer object names
		until they are first bound by calling BindBuffer.

		Error GL.INVALID_VALUE is generated if n is negative. GL.INVALID_OPERATION
		is generated if GenBuffers is executed between the execution of Begin
		and the corresponding execution of End.

		{{funcSince . "1.5+"}}
	`,
}, {
	name: "GetAttribLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	result: "glbase.Attrib",
	doc: `
		queries the previously linked program object specified
		by program for the attribute variable specified by name and returns the
		index of the generic vertex attribute that is bound to that attribute
		variable. If name is a matrix attribute variable, the index of the first
		column of the matrix is returned. If the named attribute variable is not
		an active attribute in the specified program object or if name starts with
		the reserved prefix "gl_", a value of -1 is returned.

		The association between an attribute variable name and a generic attribute
		index can be specified at any time by calling BindAttribLocation.
		Attribute bindings do not go into effect until LinkProgram is called.
		After a program object has been linked successfully, the index values for
		attribute variables remain fixed until the next link command occurs. The
		attribute values can only be queried after a link if the link was
		successful. GetAttribLocation returns the binding that actually went
		into effect the last time LinkProgram was called for the specified
		program object. Attribute bindings that have been specified since the last
		link operation are not returned by GetAttribLocation.

		Error GL_INVALID_OPERATION is generated if program is not a value
		generated by OpenGL. GL_INVALID_OPERATION is generated if program is not
		a program object. GL_INVALID_OPERATION is generated if program has not
		been successfully linked.  GL_INVALID_OPERATION is generated if
		GetAttribLocation is executed between the execution of Begin and the
		corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetShaderiv",
	doc: `
		GetShader returns in params the value of a parameter for a specific
		shader object. The following parameters are defined:

		  GL.SHADER_TYPE
		    params returns GL.VERTEX_SHADER if shader is a vertex shader object,
		    and GL.FRAGMENT_SHADER if shader is a fragment shader object.

		  GL.DELETE_STATUS
		    params returns GL.TRUE if shader is currently flagged for deletion,
		    and GL.FALSE otherwise.

		  GL.COMPILE_STATUS
		    params returns GL.TRUE if the last compile operation on shader was
		    successful, and GL.FALSE otherwise.

		  GL.INFO_LOG_LENGTH
		    params returns the number of characters in the information log for
		    shader including the null termination character (i.e., the size of the
		    character buffer required to store the information log). If shader has
		    no information log, a value of 0 is returned.

		  GL.SHADER_SOURCE_LENGTH
		    params returns the length of the concatenation of the source strings
		    that make up the shader source for the shader, including the null
		    termination character. (i.e., the size of the character buffer
		    required to store the shader source). If no source code exists, 0 is
		    returned.

		Error GL.INVALID_VALUE is generated if shader is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if shader does not refer to a
		shader object. GL.INVALID_ENUM is generated if pname is not an accepted
		value. GL.INVALID_OPERATION is generated if GetShader is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetShaderInfoLog",
	params: paramTweaks{
		"bufSize": {omit: true},
		"length":  {omit: true, single: true},
		"infoLog": {output: true, unnamed: true},
	},
	before: `
		var params [1]int32
		var length int32
		gl.GetShaderiv(shader, INFO_LOG_LENGTH, params[:])
		bufSize := params[0]
		infoLog := make([]byte, int(bufSize))
	`,
	doc: `
		returns the information log for the specified shader
		object. The information log for a shader object is modified when the
		shader is compiled.

		The information log for a shader object is a string that may contain
		diagnostic messages, warning messages, and other information about the
		last compile operation. When a shader object is created, its information
		log will be a string of length 0, and the size of the current log can be
		obtained by calling GetShaderiv with the value GL.INFO_LOG_LENGTH.

		The information log for a shader object is the OpenGL implementer's
		primary mechanism for conveying information about the compilation process.
		Therefore, the information log can be helpful to application developers
		during the development process, even when compilation is successful.
		Application developers should not expect different OpenGL implementations
		to produce identical information logs.

		Error GL.INVALID_VALUE is generated if shader is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if shader is not a shader
		object. GL.INVALID_VALUE is generated if maxLength is less than 0.
		GL.INVALID_OPERATION is generated if GetShaderInfoLog is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetUniformLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	result: "glbase.Uniform",
	doc: `
		returns an integer that represents the location of a
		specific uniform variable within a program object. name must be an active
		uniform variable name in program that is not a structure, an array of
		structures, or a subcomponent of a vector or a matrix. This function
		returns -1 if name does not correspond to an active uniform variable in
		program or if name starts with the reserved prefix "gl_".

		Uniform variables that are structures or arrays of structures may be
		queried by calling GetUniformLocation for each field within the
		structure. The array element operator "[]" and the structure field
		operator "." may be used in name in order to select elements within an
		array or fields within a structure. The result of using these operators is
		not allowed to be another structure, an array of structures, or a
		subcomponent of a vector or a matrix. Except if the last part of name
		indicates a uniform variable array, the location of the first element of
		an array can be retrieved by using the name of the array, or by using the
		name appended by "[0]".

		The actual locations assigned to uniform variables are not known until the
		program object is linked successfully. After linking has occurred, the
		command GetUniformLocation can be used to obtain the location of a
		uniform variable. This location value can then be passed to Uniform to
		set the value of the uniform variable or to GetUniform in order to query
		the current value of the uniform variable. After a program object has been
		linked successfully, the index values for uniform variables remain fixed
		until the next link command occurs. Uniform variable locations and values
		can only be queried after a link if the link was successful.

		Error GL.INVALID_VALUE is generated if program is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if program is not a program object.
		GL.INVALID_OPERATION is generated if program has not been successfully
		linked. GL.INVALID_OPERATION is generated if GetUniformLocation is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetUniformfv",
	params: paramTweaks{
		"params": {replace: true},
	},
	before: `
		var params_c [4]{{paramGoType . "params"}}
	`,
	after: `
		copy(params, params_c[:])
	`,
	doc: `
		returns in params the value of the specified uniform
		variable. The type of the uniform variable specified by location
		determines the number of values returned. If the uniform variable is
		defined in the shader as a boolean, int, or float, a single value will be
		returned. If it is defined as a vec2, ivec2, or bvec2, two values will be
		returned. If it is defined as a vec3, ivec3, or bvec3, three values will
		be returned, and so on. To query values stored in uniform variables
		declared as arrays, call {{.Name}} for each element of the array. To
		query values stored in uniform variables declared as structures, call
		{{.Name}} for each field in the structure. The values for uniform
		variables declared as a matrix will be returned in column major order.

		The locations assigned to uniform variables are not known until the
		program object is linked. After linking has occurred, the command
		GetUniformLocation can be used to obtain the location of a uniform
		variable. This location value can then be passed to {{.Name}} in order
		to query the current value of the uniform variable. After a program object
		has been linked successfully, the index values for uniform variables
		remain fixed until the next link command occurs. The uniform variable
		values can only be queried after a link if the link was successful.

		Error GL.INVALID_VALUE is generated if program is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if program is not a program
		object. GL.INVALID_OPERATION is generated if program has not been
		successfully linked. GL.INVALID_OPERATION is generated if location does
		not correspond to a valid uniform variable location for the specified
		program object. GL.INVALID_OPERATION is generated if {{.Name}} is
		executed between the execution of Begin and the corresponding execution of
		End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetUniformiv",
	copy: "GetUniformfv",
}, {
	name: "GetVertexAttribdv",
	params: paramTweaks{
		"params": {replace: true},
	},
	before: `
		var params_c [4]{{paramGoType . "params"}}
	`,
	after: `
		copy(params, params_c[:])
	`,
	doc: `
		returns in params the value of a generic vertex attribute
		parameter. The generic vertex attribute to be queried is specified by
		index, and the parameter to be queried is specified by pname.

		The accepted parameter names are as follows:

		  GL.VERTEX_ATTRIB_ARRAY_BUFFER_BINDING
		      params returns a single value, the name of the buffer object
		      currently bound to the binding point corresponding to generic vertex
		      attribute array index. If no buffer object is bound, 0 is returned.
		      The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_ENABLED
		      params returns a single value that is non-zero (true) if the vertex
		      attribute array for index is enabled and 0 (false) if it is
		      disabled. The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_SIZE
		      params returns a single value, the size of the vertex attribute
		      array for index. The size is the number of values for each element
		      of the vertex attribute array, and it will be 1, 2, 3, or 4. The
		      initial value is 4.

		  GL.VERTEX_ATTRIB_ARRAY_STRIDE
		      params returns a single value, the array stride for (number of bytes
		      between successive elements in) the vertex attribute array for
		      index. A value of 0 indicates that the array elements are stored
		      sequentially in memory. The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_TYPE
		      params returns a single value, a symbolic constant indicating the
		      array type for the vertex attribute array for index. Possible values
		      are GL.BYTE, GL.UNSIGNED_BYTE, GL.SHORT, GL.UNSIGNED_SHORT, GL.INT,
		      GL.UNSIGNED_INT, GL.FLOAT, and GL.DOUBLE. The initial value is
		      GL.FLOAT.

		  GL.VERTEX_ATTRIB_ARRAY_NORMALIZED
		      params returns a single value that is non-zero (true) if fixed-point
		      data types for the vertex attribute array indicated by index are
		      normalized when they are converted to floating point, and 0 (false)
		      otherwise. The initial value is 0.

		  GL.CURRENT_VERTEX_ATTRIB
		      params returns four values that represent the current value for the
		      generic vertex attribute specified by index. Generic vertex
		      attribute 0 is unique in that it has no current state, so an error
		      will be generated if index is 0. The initial value for all other
		      generic vertex attributes is (0,0,0,1).

		All of the parameters except GL.CURRENT_VERTEX_ATTRIB represent
		client-side state.

		Error GL.INVALID_VALUE is generated if index is greater than or equal to
		GL.MAX_VERTEX_ATTRIBS. GL.INVALID_ENUM is generated if pname is not an
		accepted value.  GL.INVALID_OPERATION is generated if index is 0 and pname
		is GL.CURRENT_VERTEX_ATTRIB.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "GetVertexAttribfv",
	copy: "GetVertexAttribdv",
}, {
	name: "GetVertexAttribiv",
	copy: "GetVertexAttribdv",
}, {
	name: "LinkProgram",
	doc: `
		links the program object specified by program. If any shader
		objects of type GL.VERTEX_SHADER are attached to program, they will be
		used to create an executable that will run on the programmable vertex
		processor. If any shader objects of type GL.FRAGMENT_SHADER are attached
		to program, they will be used to create an executable that will run on the
		programmable fragment processor.

		The status of the link operation will be stored as part of the program
		object's state. This value will be set to GL.TRUE if the program object
		was linked without errors and is ready for use, and GL.FALSE otherwise. It
		can be queried by calling GetProgramiv with arguments program and
		GL.LINK_STATUS.

		As a result of a successful link operation, all active user-defined
		uniform variables belonging to program will be initialized to 0, and each
		of the program object's active uniform variables will be assigned a
		location that can be queried by calling GetUniformLocation. Also, any
		active user-defined attribute variables that have not been bound to a
		generic vertex attribute index will be bound to one at this time.

		Linking of a program object can fail for a number of reasons as specified
		in the OpenGL Shading Language Specification. The following lists some of
		the conditions that will cause a link error.

		  - The number of active attribute variables supported by the
		    implementation has been exceeded.

		  - The storage limit for uniform variables has been exceeded.

		  - The number of active uniform variables supported by the implementation
		    has been exceeded.

		  - The main function is missing for the vertex shader or the fragment
		    shader.

		  - A varying variable actually used in the fragment shader is not
		    declared in the same way (or is not declared at all) in the vertex
		    shader.

		  - A reference to a function or variable name is unresolved.

		  - A shared global is declared with two different types or two different
		    initial values.

		  - One or more of the attached shader objects has not been successfully
		    compiled.

		  - Binding a generic attribute matrix caused some rows of the matrix to
		    fall outside the allowed maximum of GL.MAX_VERTEX_ATTRIBS.

		  - Not enough contiguous vertex attribute slots could be found to bind
		    attribute matrices.

		When a program object has been successfully linked, the program object can
		be made part of current state by calling UseProgram. Whether or not the
		link operation was successful, the program object's information log will
		be overwritten. The information log can be retrieved by calling
		GetProgramInfoLog.

		LinkProgram will also install the generated executables as part of the
		current rendering state if the link operation was successful and the
		specified program object is already currently in use as a result of a
		previous call to UseProgram. If the program object currently in use is
		relinked unsuccessfully, its link status will be set to GL.FALSE , but the
		executables and associated state will remain part of the current state
		until a subsequent call to UseProgram removes it from use. After it is
		removed from use, it cannot be made part of current state until it has
		been successfully relinked.

		If program contains shader objects of type GL.VERTEX_SHADER but does not
		contain shader objects of type GL.FRAGMENT_SHADER, the vertex shader will
		be linked against the implicit interface for fixed functionality fragment
		processing. Similarly, if program contains shader objects of type
		GL.FRAGMENT_SHADER but it does not contain shader objects of type
		GL.VERTEX_SHADER, the fragment shader will be linked against the implicit
		interface for fixed functionality vertex processing.

		The program object's information log is updated and the program is
		generated at the time of the link operation. After the link operation,
		applications are free to modify attached shader objects, compile attached
		shader objects, detach shader objects, delete shader objects, and attach
		additional shader objects. None of these operations affects the
		information log or the program that is part of the program object.

		If the link operation is unsuccessful, any information about a previous
		link operation on program is lost (i.e., a failed link does not restore
		the old state of program ). Certain information can still be retrieved
		from program even after an unsuccessful link operation. See for instance
		GetActiveAttrib and GetActiveUniform.

		Error GL.INVALID_VALUE is generated if program is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if program is not a program
		object. GL.INVALID_OPERATION is generated if LinkProgram is executed
		between the execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "MultMatrixd",
	before: `
		if len(m) != 16 {
			panic("parameter m must have length 16 for the 4x4 matrix")
		}
	`,
	doc: `
		multiplies the current matrix with the provided matrix.
		
		The m parameter must hold 16 consecutive elements of a 4x4 column-major matrix.

		The current matrix is determined by the current matrix mode (see
		MatrixMode). It is either the projection matrix, modelview matrix, or the
		texture matrix.

		For example, if the current matrix is C and the coordinates to be transformed
		are v = (v[0], v[1], v[2], v[3]), then the current transformation is C × v, or

		    c[0]  c[4]  c[8]  c[12]     v[0]
		    c[1]  c[5]  c[9]  c[13]     v[1]
		    c[2]  c[6]  c[10] c[14]  X  v[2]
		    c[3]  c[7]  c[11] c[15]     v[3]

		Calling MultMatrix with an argument of m = m[0], m[1], ..., m[15]
		replaces the current transformation with (C X M) x v, or
		
		    c[0]  c[4]  c[8]  c[12]   m[0]  m[4]  m[8]  m[12]   v[0]
		    c[1]  c[5]  c[9]  c[13]   m[1]  m[5]  m[9]  m[13]   v[1]
		    c[2]  c[6]  c[10] c[14] X m[2]  m[6]  m[10] m[14] X v[2]
		    c[3]  c[7]  c[11] c[15]   m[3]  m[7]  m[11] m[15]   v[3]

		Where 'X' denotes matrix multiplication, and v is represented as a 4x1 matrix.

		While the elements of the matrix may be specified with single or double
		precision, the GL may store or operate on these values in less-than-single
		precision.

		In many computer languages, 4×4 arrays are represented in row-major
		order. The transformations just described represent these matrices in
		column-major order. The order of the multiplication is important. For
		example, if the current transformation is a rotation, and MultMatrix is
		called with a translation matrix, the translation is done directly on the
		coordinates to be transformed, while the rotation is done on the results
		of that translation.

		GL.INVALID_OPERATION is generated if MultMatrix is executed between the
		execution of Begin and the corresponding execution of End.
	`,
}, {
	name: "MultMatrixf",
	copy: "MultMatrixd",
}, {
	name: "ShaderSource",
	params: paramTweaks{
		"glstring": {rename: "source", retype: "...string"},
		"length":   {omit: true},
		"count":    {omit: true},
	},
	before: `
		count := len(source)
		length := make([]int32, count)
		glstring := make([]unsafe.Pointer, count)
		for i, src := range source {
			length[i] = int32(len(src))
			if len(src) > 0 {
				glstring[i] = *(*unsafe.Pointer)(unsafe.Pointer(&src))
			} else {
				glstring[i] = unsafe.Pointer(uintptr(0))
			}
		}
	`,
	doc: `
		sets the source code in shader to the provided source code. Any source
		code previously stored in the shader object is completely replaced.

		Error GL.INVALID_VALUE is generated if shader is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if shader is not a shader
		object. GL.INVALID_VALUE is generated if count is less than 0.
		GL.INVALID_OPERATION is generated if ShaderSource is executed between the
		execution of Begin and the corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}, {
	name: "UseProgram",
	doc: `
		installs the program object specified by program as part of
		current rendering state. One or more executables are created in a program
		object by successfully attaching shader objects to it with AttachShader,
		successfully compiling the shader objects with CompileShader, and
		successfully linking the program object with LinkProgram.

		A program object will contain an executable that will run on the vertex
		processor if it contains one or more shader objects of type
		GL.VERTEX_SHADER that have been successfully compiled and linked.
		Similarly, a program object will contain an executable that will run on
		the fragment processor if it contains one or more shader objects of type
		GL.FRAGMENT_SHADER that have been successfully compiled and linked.

		Successfully installing an executable on a programmable processor will
		cause the corresponding fixed functionality of OpenGL to be disabled.
		Specifically, if an executable is installed on the vertex processor, the
		OpenGL fixed functionality will be disabled as follows.

		  - The modelview matrix is not applied to vertex coordinates.

		  - The projection matrix is not applied to vertex coordinates.

		  - The texture matrices are not applied to texture coordinates.

		  - Normals are not transformed to eye coordinates.

		  - Normals are not rescaled or normalized.

		  - Normalization of GL.AUTO_NORMAL evaluated normals is not performed.

		  - Texture coordinates are not generated automatically.

		  - Per-vertex lighting is not performed.

		  - Color material computations are not performed.

		  - Color index lighting is not performed.

		  - This list also applies when setting the current raster position.

		The executable that is installed on the vertex processor is expected to
		implement any or all of the desired functionality from the preceding list.
		Similarly, if an executable is installed on the fragment processor, the
		OpenGL fixed functionality will be disabled as follows.

		  - Texture environment and texture functions are not applied.

		  - Texture application is not applied.

		  - Color sum is not applied.

		  - Fog is not applied.

		Again, the fragment shader that is installed is expected to implement any
		or all of the desired functionality from the preceding list.

		While a program object is in use, applications are free to modify attached
		shader objects, compile attached shader objects, attach additional shader
		objects, and detach or delete shader objects. None of these operations
		will affect the executables that are part of the current state. However,
		relinking the program object that is currently in use will install the
		program object as part of the current rendering state if the link
		operation was successful (see LinkProgram). If the program object
		currently in use is relinked unsuccessfully, its link status will be set
		to GL.FALSE, but the executables and associated state will remain part of
		the current state until a subsequent call to UseProgram removes it from
		use. After it is removed from use, it cannot be made part of current state
		until it has been successfully relinked.

		If program contains shader objects of type GL.VERTEX_SHADER but it does
		not contain shader objects of type GL.FRAGMENT_SHADER, an executable will
		be installed on the vertex processor, but fixed functionality will be used
		for fragment processing. Similarly, if program contains shader objects of
		type GL.FRAGMENT_SHADER but it does not contain shader objects of type
		GL.VERTEX_SHADER, an executable will be installed on the fragment
		processor, but fixed functionality will be used for vertex processing. If
		program is 0, the programmable processors will be disabled, and fixed
		functionality will be used for both vertex and fragment processing.

		While a program object is in use, the state that controls the disabled
		fixed functionality may also be updated using the normal OpenGL calls.

		Like display lists and texture objects, the name space for program objects
		may be shared across a set of contexts, as long as the server sides of the
		contexts share the same address space. If the name space is shared across
		contexts, any attached objects and the data associated with those attached
		objects are shared as well.

		Applications are responsible for providing the synchronization across API
		calls when objects are accessed from different execution threads.

		Error GL.INVALID_VALUE is generated if program is neither 0 nor a value
		generated by OpenGL.  GL.INVALID_OPERATION is generated if program is not
		a program object.  GL.INVALID_OPERATION is generated if program could not
		be made part of current state.  GL.INVALID_OPERATION is generated if
		UseProgram is executed between the execution of Begin and the
		corresponding execution of End.

		{{funcSince . "2.0+"}}
	`,
}}

// vim:ts=8:tw=90:noet
