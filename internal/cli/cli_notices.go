package cli

// ===========================================================================================================
// This file was generated automatically at 15-01-2023 11:33:49 using gogenlicense.
// Do not edit manually, as changes may be overwritten.
// ===========================================================================================================

// LegalNotices contains legal and license information of external software included in this program.
// These notices consist of a list of dependencies along with their license information.
// This string is intended to be displayed to the enduser on demand.
//
// Even though the value of this variable is fixed at compile time it is omitted from this documentation.
// Instead the list of go modules, along with their licenses, is listed below.
//
// # Go Standard Library
//
// The Go Standard Library is licensed under the Terms of the BSD-3-Clause License.
// See also https://golang.org/LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//		* Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//		* Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//		* Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module gorm io gorm
//
// The Module gorm.io/gorm is licensed under the Terms of the MIT License.
// See also https://github.com/go-gorm/gorm/blob/v1.24.2/License.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in
//	all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//	THE SOFTWARE.
//
// # Module gorm io driver mysql
//
// The Module gorm.io/driver/mysql is licensed under the Terms of the MIT License.
// See also https://github.com/go-gorm/mysql/blob/v1.4.4/License.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in
//	all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//	THE SOFTWARE.
//
// # Module gopkg in yaml v2
//
// The Module gopkg.in/yaml.v2 is licensed under the Terms of the Apache-2.0 License.
// See also https://github.com/go-yaml/yaml/blob/v2.3.0/LICENSE.
//
//	                              Apache License
//	                        Version 2.0, January 2004
//	                     http://www.apache.org/licenses/
//
//	TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
//
//	1. Definitions.
//
//	   "License" shall mean the terms and conditions for use, reproduction,
//	   and distribution as defined by Sections 1 through 9 of this document.
//
//	   "Licensor" shall mean the copyright owner or entity authorized by
//	   the copyright owner that is granting the License.
//
//	   "Legal Entity" shall mean the union of the acting entity and all
//	   other entities that control, are controlled by, or are under common
//	   control with that entity. For the purposes of this definition,
//	   "control" means (i) the power, direct or indirect, to cause the
//	   direction or management of such entity, whether by contract or
//	   otherwise, or (ii) ownership of fifty percent (50%) or more of the
//	   outstanding shares, or (iii) beneficial ownership of such entity.
//
//	   "You" (or "Your") shall mean an individual or Legal Entity
//	   exercising permissions granted by this License.
//
//	   "Source" form shall mean the preferred form for making modifications,
//	   including but not limited to software source code, documentation
//	   source, and configuration files.
//
//	   "Object" form shall mean any form resulting from mechanical
//	   transformation or translation of a Source form, including but
//	   not limited to compiled object code, generated documentation,
//	   and conversions to other media types.
//
//	   "Work" shall mean the work of authorship, whether in Source or
//	   Object form, made available under the License, as indicated by a
//	   copyright notice that is included in or attached to the work
//	   (an example is provided in the Appendix below).
//
//	   "Derivative Works" shall mean any work, whether in Source or Object
//	   form, that is based on (or derived from) the Work and for which the
//	   editorial revisions, annotations, elaborations, or other modifications
//	   represent, as a whole, an original work of authorship. For the purposes
//	   of this License, Derivative Works shall not include works that remain
//	   separable from, or merely link (or bind by name) to the interfaces of,
//	   the Work and Derivative Works thereof.
//
//	   "Contribution" shall mean any work of authorship, including
//	   the original version of the Work and any modifications or additions
//	   to that Work or Derivative Works thereof, that is intentionally
//	   submitted to Licensor for inclusion in the Work by the copyright owner
//	   or by an individual or Legal Entity authorized to submit on behalf of
//	   the copyright owner. For the purposes of this definition, "submitted"
//	   means any form of electronic, verbal, or written communication sent
//	   to the Licensor or its representatives, including but not limited to
//	   communication on electronic mailing lists, source code control systems,
//	   and issue tracking systems that are managed by, or on behalf of, the
//	   Licensor for the purpose of discussing and improving the Work, but
//	   excluding communication that is conspicuously marked or otherwise
//	   designated in writing by the copyright owner as "Not a Contribution."
//
//	   "Contributor" shall mean Licensor and any individual or Legal Entity
//	   on behalf of whom a Contribution has been received by Licensor and
//	   subsequently incorporated within the Work.
//
//	2. Grant of Copyright License. Subject to the terms and conditions of
//	   this License, each Contributor hereby grants to You a perpetual,
//	   worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//	   copyright license to reproduce, prepare Derivative Works of,
//	   publicly display, publicly perform, sublicense, and distribute the
//	   Work and such Derivative Works in Source or Object form.
//
//	3. Grant of Patent License. Subject to the terms and conditions of
//	   this License, each Contributor hereby grants to You a perpetual,
//	   worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//	   (except as stated in this section) patent license to make, have made,
//	   use, offer to sell, sell, import, and otherwise transfer the Work,
//	   where such license applies only to those patent claims licensable
//	   by such Contributor that are necessarily infringed by their
//	   Contribution(s) alone or by combination of their Contribution(s)
//	   with the Work to which such Contribution(s) was submitted. If You
//	   institute patent litigation against any entity (including a
//	   cross-claim or counterclaim in a lawsuit) alleging that the Work
//	   or a Contribution incorporated within the Work constitutes direct
//	   or contributory patent infringement, then any patent licenses
//	   granted to You under this License for that Work shall terminate
//	   as of the date such litigation is filed.
//
//	4. Redistribution. You may reproduce and distribute copies of the
//	   Work or Derivative Works thereof in any medium, with or without
//	   modifications, and in Source or Object form, provided that You
//	   meet the following conditions:
//
//	   (a) You must give any other recipients of the Work or
//	       Derivative Works a copy of this License; and
//
//	   (b) You must cause any modified files to carry prominent notices
//	       stating that You changed the files; and
//
//	   (c) You must retain, in the Source form of any Derivative Works
//	       that You distribute, all copyright, patent, trademark, and
//	       attribution notices from the Source form of the Work,
//	       excluding those notices that do not pertain to any part of
//	       the Derivative Works; and
//
//	   (d) If the Work includes a "NOTICE" text file as part of its
//	       distribution, then any Derivative Works that You distribute must
//	       include a readable copy of the attribution notices contained
//	       within such NOTICE file, excluding those notices that do not
//	       pertain to any part of the Derivative Works, in at least one
//	       of the following places: within a NOTICE text file distributed
//	       as part of the Derivative Works; within the Source form or
//	       documentation, if provided along with the Derivative Works; or,
//	       within a display generated by the Derivative Works, if and
//	       wherever such third-party notices normally appear. The contents
//	       of the NOTICE file are for informational purposes only and
//	       do not modify the License. You may add Your own attribution
//	       notices within Derivative Works that You distribute, alongside
//	       or as an addendum to the NOTICE text from the Work, provided
//	       that such additional attribution notices cannot be construed
//	       as modifying the License.
//
//	   You may add Your own copyright statement to Your modifications and
//	   may provide additional or different license terms and conditions
//	   for use, reproduction, or distribution of Your modifications, or
//	   for any such Derivative Works as a whole, provided Your use,
//	   reproduction, and distribution of the Work otherwise complies with
//	   the conditions stated in this License.
//
//	5. Submission of Contributions. Unless You explicitly state otherwise,
//	   any Contribution intentionally submitted for inclusion in the Work
//	   by You to the Licensor shall be under the terms and conditions of
//	   this License, without any additional terms or conditions.
//	   Notwithstanding the above, nothing herein shall supersede or modify
//	   the terms of any separate license agreement you may have executed
//	   with Licensor regarding such Contributions.
//
//	6. Trademarks. This License does not grant permission to use the trade
//	   names, trademarks, service marks, or product names of the Licensor,
//	   except as required for reasonable and customary use in describing the
//	   origin of the Work and reproducing the content of the NOTICE file.
//
//	7. Disclaimer of Warranty. Unless required by applicable law or
//	   agreed to in writing, Licensor provides the Work (and each
//	   Contributor provides its Contributions) on an "AS IS" BASIS,
//	   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
//	   implied, including, without limitation, any warranties or conditions
//	   of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
//	   PARTICULAR PURPOSE. You are solely responsible for determining the
//	   appropriateness of using or redistributing the Work and assume any
//	   risks associated with Your exercise of permissions under this License.
//
//	8. Limitation of Liability. In no event and under no legal theory,
//	   whether in tort (including negligence), contract, or otherwise,
//	   unless required by applicable law (such as deliberate and grossly
//	   negligent acts) or agreed to in writing, shall any Contributor be
//	   liable to You for damages, including any direct, indirect, special,
//	   incidental, or consequential damages of any character arising as a
//	   result of this License or out of the use or inability to use the
//	   Work (including but not limited to damages for loss of goodwill,
//	   work stoppage, computer failure or malfunction, or any and all
//	   other commercial damages or losses), even if such Contributor
//	   has been advised of the possibility of such damages.
//
//	9. Accepting Warranty or Additional Liability. While redistributing
//	   the Work or Derivative Works thereof, You may choose to offer,
//	   and charge a fee for, acceptance of support, warranty, indemnity,
//	   or other liability obligations and/or rights consistent with this
//	   License. However, in accepting such obligations, You may act only
//	   on Your own behalf and on Your sole responsibility, not on behalf
//	   of any other Contributor, and only if You agree to indemnify,
//	   defend, and hold each Contributor harmless for any liability
//	   incurred by, or claims asserted against, such Contributor by reason
//	   of your accepting any such warranty or additional liability.
//
//	END OF TERMS AND CONDITIONS
//
//	APPENDIX: How to apply the Apache License to your work.
//
//	   To apply the Apache License to your work, attach the following
//	   boilerplate notice, with the fields enclosed by brackets "{}"
//	   replaced with your own identifying information. (Don't include
//	   the brackets!)  The text should be enclosed in the appropriate
//	   comment syntax for the file format. We also recommend that a
//	   file or class name and description of purpose be included on the
//	   same "printed page" as the copyright notice for easier
//	   identification within third-party archives.
//
//	Copyright {yyyy} {name of copyright owner}
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
//
// # Module golang org x tools go analysis
//
// The Module golang.org/x/tools/go/analysis is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/tools/+/v0.4.0:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module golang org x term
//
// The Module golang.org/x/term is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/term/+/v0.3.0:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module golang org x sys unix
//
// The Module golang.org/x/sys/unix is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/sys/+/v0.3.0:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module golang org x sync errgroup
//
// The Module golang.org/x/sync/errgroup is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/sync/+/v0.1.0:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module golang org x exp
//
// The Module golang.org/x/exp is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/exp/+/8879d019:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module golang org x crypto
//
// The Module golang.org/x/crypto is licensed under the Terms of the BSD-3-Clause License.
// See also https://cs.opensource.google/go/x/crypto/+/v0.3.0:LICENSE.
//
//	Copyright (c) 2009 The Go Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com yuin goldmark meta
//
// The Module github.com/yuin/goldmark-meta is licensed under the Terms of the MIT License.
// See also https://github.com/yuin/goldmark-meta/blob/v1.1.0/LICENSE.
//
//	MIT License
//
//	Copyright (c) 2019 Yusuke Inuzuka
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com yuin goldmark
//
// The Module github.com/yuin/goldmark is licensed under the Terms of the MIT License.
// See also https://github.com/yuin/goldmark/blob/v1.4.13/LICENSE.
//
//	MIT License
//
//	Copyright (c) 2019 Yusuke Inuzuka
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com tkw1536 goprogram
//
// The Module github.com/tkw1536/goprogram is licensed under the Terms of the MIT License.
// See also https://github.com/tkw1536/goprogram/blob/v0.2.4/LICENSE.
//
//	MIT License
//
//	Copyright (c) 2022 Tom Wiesing
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com tdewolff parse
//
// The Module github.com/tdewolff/parse is licensed under the Terms of the MIT License.
// See also https://github.com/tdewolff/parse/blob/v2.3.4/LICENSE.md.
//
//	Copyright (c) 2015 Taco de Wolff
//
//	 Permission is hereby granted, free of charge, to any person
//	 obtaining a copy of this software and associated documentation
//	 files (the "Software"), to deal in the Software without
//	 restriction, including without limitation the rights to use,
//	 copy, modify, merge, publish, distribute, sublicense, and/or sell
//	 copies of the Software, and to permit persons to whom the
//	 Software is furnished to do so, subject to the following
//	 conditions:
//
//	 The above copyright notice and this permission notice shall be
//	 included in all copies or substantial portions of the Software.
//
//	 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//	 EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
//	 OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//	 NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//	 HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
//	 WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
//	 FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//	 OTHER DEALINGS IN THE SOFTWARE.
//
// # Module github com tdewolff minify
//
// The Module github.com/tdewolff/minify is licensed under the Terms of the MIT License.
// See also https://github.com/tdewolff/minify/blob/v2.3.6/LICENSE.md.
//
//	Copyright (c) 2015 Taco de Wolff
//
//	 Permission is hereby granted, free of charge, to any person
//	 obtaining a copy of this software and associated documentation
//	 files (the "Software"), to deal in the Software without
//	 restriction, including without limitation the rights to use,
//	 copy, modify, merge, publish, distribute, sublicense, and/or sell
//	 copies of the Software, and to permit persons to whom the
//	 Software is furnished to do so, subject to the following
//	 conditions:
//
//	 The above copyright notice and this permission notice shall be
//	 included in all copies or substantial portions of the Software.
//
//	 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//	 EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
//	 OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//	 NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//	 HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
//	 WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
//	 FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//	 OTHER DEALINGS IN THE SOFTWARE.
//
// # Module github com rs zerolog
//
// The Module github.com/rs/zerolog is licensed under the Terms of the MIT License.
// See also https://github.com/rs/zerolog/blob/v1.28.0/LICENSE.
//
//	MIT License
//
//	Copyright (c) 2017 Olivier Poitrey
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com pquerna otp
//
// The Module github.com/pquerna/otp is licensed under the Terms of the Apache-2.0 License.
// See also https://github.com/pquerna/otp/blob/v1.4.0/LICENSE.
//
//	                              Apache License
//	                        Version 2.0, January 2004
//	                     http://www.apache.org/licenses/
//
//	TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
//
//	1. Definitions.
//
//	   "License" shall mean the terms and conditions for use, reproduction,
//	   and distribution as defined by Sections 1 through 9 of this document.
//
//	   "Licensor" shall mean the copyright owner or entity authorized by
//	   the copyright owner that is granting the License.
//
//	   "Legal Entity" shall mean the union of the acting entity and all
//	   other entities that control, are controlled by, or are under common
//	   control with that entity. For the purposes of this definition,
//	   "control" means (i) the power, direct or indirect, to cause the
//	   direction or management of such entity, whether by contract or
//	   otherwise, or (ii) ownership of fifty percent (50%) or more of the
//	   outstanding shares, or (iii) beneficial ownership of such entity.
//
//	   "You" (or "Your") shall mean an individual or Legal Entity
//	   exercising permissions granted by this License.
//
//	   "Source" form shall mean the preferred form for making modifications,
//	   including but not limited to software source code, documentation
//	   source, and configuration files.
//
//	   "Object" form shall mean any form resulting from mechanical
//	   transformation or translation of a Source form, including but
//	   not limited to compiled object code, generated documentation,
//	   and conversions to other media types.
//
//	   "Work" shall mean the work of authorship, whether in Source or
//	   Object form, made available under the License, as indicated by a
//	   copyright notice that is included in or attached to the work
//	   (an example is provided in the Appendix below).
//
//	   "Derivative Works" shall mean any work, whether in Source or Object
//	   form, that is based on (or derived from) the Work and for which the
//	   editorial revisions, annotations, elaborations, or other modifications
//	   represent, as a whole, an original work of authorship. For the purposes
//	   of this License, Derivative Works shall not include works that remain
//	   separable from, or merely link (or bind by name) to the interfaces of,
//	   the Work and Derivative Works thereof.
//
//	   "Contribution" shall mean any work of authorship, including
//	   the original version of the Work and any modifications or additions
//	   to that Work or Derivative Works thereof, that is intentionally
//	   submitted to Licensor for inclusion in the Work by the copyright owner
//	   or by an individual or Legal Entity authorized to submit on behalf of
//	   the copyright owner. For the purposes of this definition, "submitted"
//	   means any form of electronic, verbal, or written communication sent
//	   to the Licensor or its representatives, including but not limited to
//	   communication on electronic mailing lists, source code control systems,
//	   and issue tracking systems that are managed by, or on behalf of, the
//	   Licensor for the purpose of discussing and improving the Work, but
//	   excluding communication that is conspicuously marked or otherwise
//	   designated in writing by the copyright owner as "Not a Contribution."
//
//	   "Contributor" shall mean Licensor and any individual or Legal Entity
//	   on behalf of whom a Contribution has been received by Licensor and
//	   subsequently incorporated within the Work.
//
//	2. Grant of Copyright License. Subject to the terms and conditions of
//	   this License, each Contributor hereby grants to You a perpetual,
//	   worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//	   copyright license to reproduce, prepare Derivative Works of,
//	   publicly display, publicly perform, sublicense, and distribute the
//	   Work and such Derivative Works in Source or Object form.
//
//	3. Grant of Patent License. Subject to the terms and conditions of
//	   this License, each Contributor hereby grants to You a perpetual,
//	   worldwide, non-exclusive, no-charge, royalty-free, irrevocable
//	   (except as stated in this section) patent license to make, have made,
//	   use, offer to sell, sell, import, and otherwise transfer the Work,
//	   where such license applies only to those patent claims licensable
//	   by such Contributor that are necessarily infringed by their
//	   Contribution(s) alone or by combination of their Contribution(s)
//	   with the Work to which such Contribution(s) was submitted. If You
//	   institute patent litigation against any entity (including a
//	   cross-claim or counterclaim in a lawsuit) alleging that the Work
//	   or a Contribution incorporated within the Work constitutes direct
//	   or contributory patent infringement, then any patent licenses
//	   granted to You under this License for that Work shall terminate
//	   as of the date such litigation is filed.
//
//	4. Redistribution. You may reproduce and distribute copies of the
//	   Work or Derivative Works thereof in any medium, with or without
//	   modifications, and in Source or Object form, provided that You
//	   meet the following conditions:
//
//	   (a) You must give any other recipients of the Work or
//	       Derivative Works a copy of this License; and
//
//	   (b) You must cause any modified files to carry prominent notices
//	       stating that You changed the files; and
//
//	   (c) You must retain, in the Source form of any Derivative Works
//	       that You distribute, all copyright, patent, trademark, and
//	       attribution notices from the Source form of the Work,
//	       excluding those notices that do not pertain to any part of
//	       the Derivative Works; and
//
//	   (d) If the Work includes a "NOTICE" text file as part of its
//	       distribution, then any Derivative Works that You distribute must
//	       include a readable copy of the attribution notices contained
//	       within such NOTICE file, excluding those notices that do not
//	       pertain to any part of the Derivative Works, in at least one
//	       of the following places: within a NOTICE text file distributed
//	       as part of the Derivative Works; within the Source form or
//	       documentation, if provided along with the Derivative Works; or,
//	       within a display generated by the Derivative Works, if and
//	       wherever such third-party notices normally appear. The contents
//	       of the NOTICE file are for informational purposes only and
//	       do not modify the License. You may add Your own attribution
//	       notices within Derivative Works that You distribute, alongside
//	       or as an addendum to the NOTICE text from the Work, provided
//	       that such additional attribution notices cannot be construed
//	       as modifying the License.
//
//	   You may add Your own copyright statement to Your modifications and
//	   may provide additional or different license terms and conditions
//	   for use, reproduction, or distribution of Your modifications, or
//	   for any such Derivative Works as a whole, provided Your use,
//	   reproduction, and distribution of the Work otherwise complies with
//	   the conditions stated in this License.
//
//	5. Submission of Contributions. Unless You explicitly state otherwise,
//	   any Contribution intentionally submitted for inclusion in the Work
//	   by You to the Licensor shall be under the terms and conditions of
//	   this License, without any additional terms or conditions.
//	   Notwithstanding the above, nothing herein shall supersede or modify
//	   the terms of any separate license agreement you may have executed
//	   with Licensor regarding such Contributions.
//
//	6. Trademarks. This License does not grant permission to use the trade
//	   names, trademarks, service marks, or product names of the Licensor,
//	   except as required for reasonable and customary use in describing the
//	   origin of the Work and reproducing the content of the NOTICE file.
//
//	7. Disclaimer of Warranty. Unless required by applicable law or
//	   agreed to in writing, Licensor provides the Work (and each
//	   Contributor provides its Contributions) on an "AS IS" BASIS,
//	   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
//	   implied, including, without limitation, any warranties or conditions
//	   of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
//	   PARTICULAR PURPOSE. You are solely responsible for determining the
//	   appropriateness of using or redistributing the Work and assume any
//	   risks associated with Your exercise of permissions under this License.
//
//	8. Limitation of Liability. In no event and under no legal theory,
//	   whether in tort (including negligence), contract, or otherwise,
//	   unless required by applicable law (such as deliberate and grossly
//	   negligent acts) or agreed to in writing, shall any Contributor be
//	   liable to You for damages, including any direct, indirect, special,
//	   incidental, or consequential damages of any character arising as a
//	   result of this License or out of the use or inability to use the
//	   Work (including but not limited to damages for loss of goodwill,
//	   work stoppage, computer failure or malfunction, or any and all
//	   other commercial damages or losses), even if such Contributor
//	   has been advised of the possibility of such damages.
//
//	9. Accepting Warranty or Additional Liability. While redistributing
//	   the Work or Derivative Works thereof, You may choose to offer,
//	   and charge a fee for, acceptance of support, warranty, indemnity,
//	   or other liability obligations and/or rights consistent with this
//	   License. However, in accepting such obligations, You may act only
//	   on Your own behalf and on Your sole responsibility, not on behalf
//	   of any other Contributor, and only if You agree to indemnify,
//	   defend, and hold each Contributor harmless for any liability
//	   incurred by, or claims asserted against, such Contributor by reason
//	   of your accepting any such warranty or additional liability.
//
//	END OF TERMS AND CONDITIONS
//
//	APPENDIX: How to apply the Apache License to your work.
//
//	   To apply the Apache License to your work, attach the following
//	   boilerplate notice, with the fields enclosed by brackets "[]"
//	   replaced with your own identifying information. (Don't include
//	   the brackets!)  The text should be enclosed in the appropriate
//	   comment syntax for the file format. We also recommend that a
//	   file or class name and description of purpose be included on the
//	   same "printed page" as the copyright notice for easier
//	   identification within third-party archives.
//
//	Copyright [yyyy] [name of copyright owner]
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
//
// # Module github com pkg errors
//
// The Module github.com/pkg/errors is licensed under the Terms of the BSD-2-Clause License.
// See also https://github.com/pkg/errors/blob/v0.9.1/LICENSE.
//
//	Copyright (c) 2015, Dave Cheney <dave@cheney.net>
//	All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are met:
//
//	* Redistributions of source code must retain the above copyright notice, this
//	  list of conditions and the following disclaimer.
//
//	* Redistributions in binary form must reproduce the above copyright notice,
//	  this list of conditions and the following disclaimer in the documentation
//	  and/or other materials provided with the distribution.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
//	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
//	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
//	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
//	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
//	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
//	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
//	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com mattn go isatty
//
// The Module github.com/mattn/go-isatty is licensed under the Terms of the MIT License.
// See also https://github.com/mattn/go-isatty/blob/v0.0.16/LICENSE.
//
//	Copyright (c) Yasuhiro MATSUMOTO <mattn.jp@gmail.com>
//
//	MIT License (Expat)
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// # Module github com mattn go colorable
//
// The Module github.com/mattn/go-colorable is licensed under the Terms of the MIT License.
// See also https://github.com/mattn/go-colorable/blob/v0.1.13/LICENSE.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2016 Yasuhiro Matsumoto
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com julienschmidt httprouter
//
// The Module github.com/julienschmidt/httprouter is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/julienschmidt/httprouter/blob/v1.3.0/LICENSE.
//
//	BSD 3-Clause License
//
//	Copyright (c) 2013, Julien Schmidt
//	All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are met:
//
//	1. Redistributions of source code must retain the above copyright notice, this
//	   list of conditions and the following disclaimer.
//
//	2. Redistributions in binary form must reproduce the above copyright notice,
//	   this list of conditions and the following disclaimer in the documentation
//	   and/or other materials provided with the distribution.
//
//	3. Neither the name of the copyright holder nor the names of its
//	   contributors may be used to endorse or promote products derived from
//	   this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
//	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
//	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
//	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
//	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
//	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
//	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
//	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com jinzhu now
//
// The Module github.com/jinzhu/now is licensed under the Terms of the MIT License.
// See also https://github.com/jinzhu/now/blob/v1.1.5/License.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in
//	all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//	THE SOFTWARE.
//
// # Module github com jinzhu inflection
//
// The Module github.com/jinzhu/inflection is licensed under the Terms of the MIT License.
// See also https://github.com/jinzhu/inflection/blob/v1.0.0/LICENSE.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2015 - Jinzhu
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com jessevdk go flags
//
// The Module github.com/jessevdk/go-flags is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/jessevdk/go-flags/blob/v1.5.0/LICENSE.
//
//	Copyright (c) 2012 Jesse van den Kieboom. All rights reserved.
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	     notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	     copyright notice, this list of conditions and the following disclaimer
//	     in the documentation and/or other materials provided with the
//	     distribution.
//	   * Neither the name of Google Inc. nor the names of its
//	     contributors may be used to endorse or promote products derived from
//	     this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com gosuri uilive
//
// The Module github.com/gosuri/uilive is licensed under the Terms of the MIT License.
// See also https://github.com/gosuri/uilive/blob/v0.0.4/LICENSE.
//
//	MIT License
//	===========
//
//	Copyright (c) 2015, Greg Osuri
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// # Module github com gorilla websocket
//
// The Module github.com/gorilla/websocket is licensed under the Terms of the BSD-2-Clause License.
// See also https://github.com/gorilla/websocket/blob/v1.5.0/LICENSE.
//
//	Copyright (c) 2013 The Gorilla WebSocket Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are met:
//
//	  Redistributions of source code must retain the above copyright notice, this
//	  list of conditions and the following disclaimer.
//
//	  Redistributions in binary form must reproduce the above copyright notice,
//	  this list of conditions and the following disclaimer in the documentation
//	  and/or other materials provided with the distribution.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
//	ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
//	WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
//	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
//	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
//	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
//	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
//	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com gorilla sessions
//
// The Module github.com/gorilla/sessions is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/gorilla/sessions/blob/v1.2.1/LICENSE.
//
//	Copyright (c) 2012-2018 The Gorilla Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//		 * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//		 * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//		 * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com gorilla securecookie
//
// The Module github.com/gorilla/securecookie is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/gorilla/securecookie/blob/v1.1.1/LICENSE.
//
//	Copyright (c) 2012 Rodrigo Moraes. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//		 * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//		 * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//		 * Neither the name of Google Inc. nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com gorilla csrf
//
// The Module github.com/gorilla/csrf is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/gorilla/csrf/blob/v1.7.1/LICENSE.
//
//	Copyright (c) 2015-2018, The Gorilla Authors. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without modification,
//	are permitted provided that the following conditions are met:
//
//	1. Redistributions of source code must retain the above copyright notice, this
//	list of conditions and the following disclaimer.
//
//	2. Redistributions in binary form must reproduce the above copyright notice,
//	this list of conditions and the following disclaimer in the documentation and/or
//	other materials provided with the distribution.
//
//	3. Neither the name of the copyright holder nor the names of its contributors
//	may be used to endorse or promote products derived from this software without
//	specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
//	ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
//	WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
//	ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
//	(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
//	LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
//	ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
//	SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com go sql driver mysql
//
// The Module github.com/go-sql-driver/mysql is licensed under the Terms of the MPL-2.0 License.
// See also https://github.com/go-sql-driver/mysql/blob/v1.6.0/LICENSE.
//
//	Mozilla Public License Version 2.0
//	==================================
//
//	1. Definitions
//	--------------
//
//	1.1. "Contributor"
//	    means each individual or legal entity that creates, contributes to
//	    the creation of, or owns Covered Software.
//
//	1.2. "Contributor Version"
//	    means the combination of the Contributions of others (if any) used
//	    by a Contributor and that particular Contributor's Contribution.
//
//	1.3. "Contribution"
//	    means Covered Software of a particular Contributor.
//
//	1.4. "Covered Software"
//	    means Source Code Form to which the initial Contributor has attached
//	    the notice in Exhibit A, the Executable Form of such Source Code
//	    Form, and Modifications of such Source Code Form, in each case
//	    including portions thereof.
//
//	1.5. "Incompatible With Secondary Licenses"
//	    means
//
//	    (a) that the initial Contributor has attached the notice described
//	        in Exhibit B to the Covered Software; or
//
//	    (b) that the Covered Software was made available under the terms of
//	        version 1.1 or earlier of the License, but not also under the
//	        terms of a Secondary License.
//
//	1.6. "Executable Form"
//	    means any form of the work other than Source Code Form.
//
//	1.7. "Larger Work"
//	    means a work that combines Covered Software with other material, in
//	    a separate file or files, that is not Covered Software.
//
//	1.8. "License"
//	    means this document.
//
//	1.9. "Licensable"
//	    means having the right to grant, to the maximum extent possible,
//	    whether at the time of the initial grant or subsequently, any and
//	    all of the rights conveyed by this License.
//
//	1.10. "Modifications"
//	    means any of the following:
//
//	    (a) any file in Source Code Form that results from an addition to,
//	        deletion from, or modification of the contents of Covered
//	        Software; or
//
//	    (b) any new file in Source Code Form that contains any Covered
//	        Software.
//
//	1.11. "Patent Claims" of a Contributor
//	    means any patent claim(s), including without limitation, method,
//	    process, and apparatus claims, in any patent Licensable by such
//	    Contributor that would be infringed, but for the grant of the
//	    License, by the making, using, selling, offering for sale, having
//	    made, import, or transfer of either its Contributions or its
//	    Contributor Version.
//
//	1.12. "Secondary License"
//	    means either the GNU General Public License, Version 2.0, the GNU
//	    Lesser General Public License, Version 2.1, the GNU Affero General
//	    Public License, Version 3.0, or any later versions of those
//	    licenses.
//
//	1.13. "Source Code Form"
//	    means the form of the work preferred for making modifications.
//
//	1.14. "You" (or "Your")
//	    means an individual or a legal entity exercising rights under this
//	    License. For legal entities, "You" includes any entity that
//	    controls, is controlled by, or is under common control with You. For
//	    purposes of this definition, "control" means (a) the power, direct
//	    or indirect, to cause the direction or management of such entity,
//	    whether by contract or otherwise, or (b) ownership of more than
//	    fifty percent (50%) of the outstanding shares or beneficial
//	    ownership of such entity.
//
//	2. License Grants and Conditions
//	--------------------------------
//
//	2.1. Grants
//
//	Each Contributor hereby grants You a world-wide, royalty-free,
//	non-exclusive license:
//
//	(a) under intellectual property rights (other than patent or trademark)
//	    Licensable by such Contributor to use, reproduce, make available,
//	    modify, display, perform, distribute, and otherwise exploit its
//	    Contributions, either on an unmodified basis, with Modifications, or
//	    as part of a Larger Work; and
//
//	(b) under Patent Claims of such Contributor to make, use, sell, offer
//	    for sale, have made, import, and otherwise transfer either its
//	    Contributions or its Contributor Version.
//
//	2.2. Effective Date
//
//	The licenses granted in Section 2.1 with respect to any Contribution
//	become effective for each Contribution on the date the Contributor first
//	distributes such Contribution.
//
//	2.3. Limitations on Grant Scope
//
//	The licenses granted in this Section 2 are the only rights granted under
//	this License. No additional rights or licenses will be implied from the
//	distribution or licensing of Covered Software under this License.
//	Notwithstanding Section 2.1(b) above, no patent license is granted by a
//	Contributor:
//
//	(a) for any code that a Contributor has removed from Covered Software;
//	    or
//
//	(b) for infringements caused by: (i) Your and any other third party's
//	    modifications of Covered Software, or (ii) the combination of its
//	    Contributions with other software (except as part of its Contributor
//	    Version); or
//
//	(c) under Patent Claims infringed by Covered Software in the absence of
//	    its Contributions.
//
//	This License does not grant any rights in the trademarks, service marks,
//	or logos of any Contributor (except as may be necessary to comply with
//	the notice requirements in Section 3.4).
//
//	2.4. Subsequent Licenses
//
//	No Contributor makes additional grants as a result of Your choice to
//	distribute the Covered Software under a subsequent version of this
//	License (see Section 10.2) or under the terms of a Secondary License (if
//	permitted under the terms of Section 3.3).
//
//	2.5. Representation
//
//	Each Contributor represents that the Contributor believes its
//	Contributions are its original creation(s) or it has sufficient rights
//	to grant the rights to its Contributions conveyed by this License.
//
//	2.6. Fair Use
//
//	This License is not intended to limit any rights You have under
//	applicable copyright doctrines of fair use, fair dealing, or other
//	equivalents.
//
//	2.7. Conditions
//
//	Sections 3.1, 3.2, 3.3, and 3.4 are conditions of the licenses granted
//	in Section 2.1.
//
//	3. Responsibilities
//	-------------------
//
//	3.1. Distribution of Source Form
//
//	All distribution of Covered Software in Source Code Form, including any
//	Modifications that You create or to which You contribute, must be under
//	the terms of this License. You must inform recipients that the Source
//	Code Form of the Covered Software is governed by the terms of this
//	License, and how they can obtain a copy of this License. You may not
//	attempt to alter or restrict the recipients' rights in the Source Code
//	Form.
//
//	3.2. Distribution of Executable Form
//
//	If You distribute Covered Software in Executable Form then:
//
//	(a) such Covered Software must also be made available in Source Code
//	    Form, as described in Section 3.1, and You must inform recipients of
//	    the Executable Form how they can obtain a copy of such Source Code
//	    Form by reasonable means in a timely manner, at a charge no more
//	    than the cost of distribution to the recipient; and
//
//	(b) You may distribute such Executable Form under the terms of this
//	    License, or sublicense it under different terms, provided that the
//	    license for the Executable Form does not attempt to limit or alter
//	    the recipients' rights in the Source Code Form under this License.
//
//	3.3. Distribution of a Larger Work
//
//	You may create and distribute a Larger Work under terms of Your choice,
//	provided that You also comply with the requirements of this License for
//	the Covered Software. If the Larger Work is a combination of Covered
//	Software with a work governed by one or more Secondary Licenses, and the
//	Covered Software is not Incompatible With Secondary Licenses, this
//	License permits You to additionally distribute such Covered Software
//	under the terms of such Secondary License(s), so that the recipient of
//	the Larger Work may, at their option, further distribute the Covered
//	Software under the terms of either this License or such Secondary
//	License(s).
//
//	3.4. Notices
//
//	You may not remove or alter the substance of any license notices
//	(including copyright notices, patent notices, disclaimers of warranty,
//	or limitations of liability) contained within the Source Code Form of
//	the Covered Software, except that You may alter any license notices to
//	the extent required to remedy known factual inaccuracies.
//
//	3.5. Application of Additional Terms
//
//	You may choose to offer, and to charge a fee for, warranty, support,
//	indemnity or liability obligations to one or more recipients of Covered
//	Software. However, You may do so only on Your own behalf, and not on
//	behalf of any Contributor. You must make it absolutely clear that any
//	such warranty, support, indemnity, or liability obligation is offered by
//	You alone, and You hereby agree to indemnify every Contributor for any
//	liability incurred by such Contributor as a result of warranty, support,
//	indemnity or liability terms You offer. You may include additional
//	disclaimers of warranty and limitations of liability specific to any
//	jurisdiction.
//
//	4. Inability to Comply Due to Statute or Regulation
//	---------------------------------------------------
//
//	If it is impossible for You to comply with any of the terms of this
//	License with respect to some or all of the Covered Software due to
//	statute, judicial order, or regulation then You must: (a) comply with
//	the terms of this License to the maximum extent possible; and (b)
//	describe the limitations and the code they affect. Such description must
//	be placed in a text file included with all distributions of the Covered
//	Software under this License. Except to the extent prohibited by statute
//	or regulation, such description must be sufficiently detailed for a
//	recipient of ordinary skill to be able to understand it.
//
//	5. Termination
//	--------------
//
//	5.1. The rights granted under this License will terminate automatically
//	if You fail to comply with any of its terms. However, if You become
//	compliant, then the rights granted under this License from a particular
//	Contributor are reinstated (a) provisionally, unless and until such
//	Contributor explicitly and finally terminates Your grants, and (b) on an
//	ongoing basis, if such Contributor fails to notify You of the
//	non-compliance by some reasonable means prior to 60 days after You have
//	come back into compliance. Moreover, Your grants from a particular
//	Contributor are reinstated on an ongoing basis if such Contributor
//	notifies You of the non-compliance by some reasonable means, this is the
//	first time You have received notice of non-compliance with this License
//	from such Contributor, and You become compliant prior to 30 days after
//	Your receipt of the notice.
//
//	5.2. If You initiate litigation against any entity by asserting a patent
//	infringement claim (excluding declaratory judgment actions,
//	counter-claims, and cross-claims) alleging that a Contributor Version
//	directly or indirectly infringes any patent, then the rights granted to
//	You by any and all Contributors for the Covered Software under Section
//	2.1 of this License shall terminate.
//
//	5.3. In the event of termination under Sections 5.1 or 5.2 above, all
//	end user license agreements (excluding distributors and resellers) which
//	have been validly granted by You or Your distributors under this License
//	prior to termination shall survive termination.
//
//	************************************************************************
//	*                                                                      *
//	*  6. Disclaimer of Warranty                                           *
//	*  -------------------------                                           *
//	*                                                                      *
//	*  Covered Software is provided under this License on an "as is"       *
//	*  basis, without warranty of any kind, either expressed, implied, or  *
//	*  statutory, including, without limitation, warranties that the       *
//	*  Covered Software is free of defects, merchantable, fit for a        *
//	*  particular purpose or non-infringing. The entire risk as to the     *
//	*  quality and performance of the Covered Software is with You.        *
//	*  Should any Covered Software prove defective in any respect, You     *
//	*  (not any Contributor) assume the cost of any necessary servicing,   *
//	*  repair, or correction. This disclaimer of warranty constitutes an   *
//	*  essential part of this License. No use of any Covered Software is   *
//	*  authorized under this License except under this disclaimer.         *
//	*                                                                      *
//	************************************************************************
//
//	************************************************************************
//	*                                                                      *
//	*  7. Limitation of Liability                                          *
//	*  --------------------------                                          *
//	*                                                                      *
//	*  Under no circumstances and under no legal theory, whether tort      *
//	*  (including negligence), contract, or otherwise, shall any           *
//	*  Contributor, or anyone who distributes Covered Software as          *
//	*  permitted above, be liable to You for any direct, indirect,         *
//	*  special, incidental, or consequential damages of any character      *
//	*  including, without limitation, damages for lost profits, loss of    *
//	*  goodwill, work stoppage, computer failure or malfunction, or any    *
//	*  and all other commercial damages or losses, even if such party      *
//	*  shall have been informed of the possibility of such damages. This   *
//	*  limitation of liability shall not apply to liability for death or   *
//	*  personal injury resulting from such party's negligence to the       *
//	*  extent applicable law prohibits such limitation. Some               *
//	*  jurisdictions do not allow the exclusion or limitation of           *
//	*  incidental or consequential damages, so this exclusion and          *
//	*  limitation may not apply to You.                                    *
//	*                                                                      *
//	************************************************************************
//
//	8. Litigation
//	-------------
//
//	Any litigation relating to this License may be brought only in the
//	courts of a jurisdiction where the defendant maintains its principal
//	place of business and such litigation shall be governed by laws of that
//	jurisdiction, without reference to its conflict-of-law provisions.
//	Nothing in this Section shall prevent a party's ability to bring
//	cross-claims or counter-claims.
//
//	9. Miscellaneous
//	----------------
//
//	This License represents the complete agreement concerning the subject
//	matter hereof. If any provision of this License is held to be
//	unenforceable, such provision shall be reformed only to the extent
//	necessary to make it enforceable. Any law or regulation which provides
//	that the language of a contract shall be construed against the drafter
//	shall not be used to construe this License against a Contributor.
//
//	10. Versions of the License
//	---------------------------
//
//	10.1. New Versions
//
//	Mozilla Foundation is the license steward. Except as provided in Section
//	10.3, no one other than the license steward has the right to modify or
//	publish new versions of this License. Each version will be given a
//	distinguishing version number.
//
//	10.2. Effect of New Versions
//
//	You may distribute the Covered Software under the terms of the version
//	of the License under which You originally received the Covered Software,
//	or under the terms of any subsequent version published by the license
//	steward.
//
//	10.3. Modified Versions
//
//	If you create software not governed by this License, and you want to
//	create a new license for such software, you may create and use a
//	modified version of this License if you rename the license and remove
//	any references to the name of the license steward (except to note that
//	such modified license differs from this License).
//
//	10.4. Distributing Source Code Form that is Incompatible With Secondary
//	Licenses
//
//	If You choose to distribute Source Code Form that is Incompatible With
//	Secondary Licenses under the terms of this version of the License, the
//	notice described in Exhibit B of this License must be attached.
//
//	Exhibit A - Source Code Form License Notice
//	-------------------------------------------
//
//	  This Source Code Form is subject to the terms of the Mozilla Public
//	  License, v. 2.0. If a copy of the MPL was not distributed with this
//	  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
//	If it is not possible or desirable to put the notice in a particular
//	file, then You may include the notice in a location (such as a LICENSE
//	file in a relevant directory) where a recipient would be likely to look
//	for such a notice.
//
//	You may add additional accurate notices of copyright ownership.
//
//	Exhibit B - "Incompatible With Secondary Licenses" Notice
//	---------------------------------------------------------
//
//	  This Source Code Form is "Incompatible With Secondary Licenses", as
//	  defined by the Mozilla Public License, v. 2.0.
//
// # Module github com gliderlabs ssh
//
// The Module github.com/gliderlabs/ssh is licensed under the Terms of the BSD-3-Clause License.
// See also https://github.com/gliderlabs/ssh/blob/v0.3.5/LICENSE.
//
//	Copyright (c) 2016 Glider Labs. All rights reserved.
//
//	Redistribution and use in source and binary forms, with or without
//	modification, are permitted provided that the following conditions are
//	met:
//
//	   * Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	   * Redistributions in binary form must reproduce the above
//	copyright notice, this list of conditions and the following disclaimer
//	in the documentation and/or other materials provided with the
//	distribution.
//	   * Neither the name of Glider Labs nor the names of its
//	contributors may be used to endorse or promote products derived from
//	this software without specific prior written permission.
//
//	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// # Module github com feiin sqlstring
//
// The Module github.com/feiin/sqlstring is licensed under the Terms of the MIT License.
// See also https://github.com/feiin/sqlstring/blob/v0.3.0/LICENSE.
//
//	MIT License
//
//	Copyright (c) 2021 solar
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com boombuler barcode
//
// The Module github.com/boombuler/barcode is licensed under the Terms of the MIT License.
// See also https://github.com/boombuler/barcode/blob/6c824513bacc/LICENSE.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2014 Florian Sundermann
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com anmitsu go shlex
//
// The Module github.com/anmitsu/go-shlex is licensed under the Terms of the MIT License.
// See also https://github.com/anmitsu/go-shlex/blob/38f4b401e2be/LICENSE.
//
//	Copyright (c) anmitsu <anmitsu.s@gmail.com>
//
//	Permission is hereby granted, free of charge, to any person obtaining
//	a copy of this software and associated documentation files (the
//	"Software"), to deal in the Software without restriction, including
//	without limitation the rights to use, copy, modify, merge, publish,
//	distribute, sublicense, and/or sell copies of the Software, and to
//	permit persons to whom the Software is furnished to do so, subject to
//	the following conditions:
//
//	The above copyright notice and this permission notice shall be
//	included in all copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//	EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//	MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//	NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
//	LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
//	OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
//	WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// # Module github com alessio shellescape
//
// The Module github.com/alessio/shellescape is licensed under the Terms of the MIT License.
// See also https://github.com/alessio/shellescape/blob/v1.4.1/LICENSE.
//
//	The MIT License (MIT)
//
//	Copyright (c) 2016 Alessio Treglia
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE.
//
// # Module github com Showmax go fqdn
//
// The Module github.com/Showmax/go-fqdn is licensed under the Terms of the Apache-2.0 License.
// See also https://github.com/Showmax/go-fqdn/blob/v1.0.0/LICENSE.
//
//	Copyright since 2015 Showmax s.r.o.
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
//
// # Module github com FAU CDI wdresolve
//
// The Module github.com/FAU-CDI/wdresolve is licensed under the Terms of the AGPL-3.0 License.
// See also https://github.com/FAU-CDI/wdresolve/blob/c9c6779d7c41/LICENSE.
//
//	                    GNU AFFERO GENERAL PUBLIC LICENSE
//	                       Version 3, 19 November 2007
//
//	 Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
//	 Everyone is permitted to copy and distribute verbatim copies
//	 of this license document, but changing it is not allowed.
//
//	                            Preamble
//
//	  The GNU Affero General Public License is a free, copyleft license for
//	software and other kinds of works, specifically designed to ensure
//	cooperation with the community in the case of network server software.
//
//	  The licenses for most software and other practical works are designed
//	to take away your freedom to share and change the works.  By contrast,
//	our General Public Licenses are intended to guarantee your freedom to
//	share and change all versions of a program--to make sure it remains free
//	software for all its users.
//
//	  When we speak of free software, we are referring to freedom, not
//	price.  Our General Public Licenses are designed to make sure that you
//	have the freedom to distribute copies of free software (and charge for
//	them if you wish), that you receive source code or can get it if you
//	want it, that you can change the software or use pieces of it in new
//	free programs, and that you know you can do these things.
//
//	  Developers that use our General Public Licenses protect your rights
//	with two steps: (1) assert copyright on the software, and (2) offer
//	you this License which gives you legal permission to copy, distribute
//	and/or modify the software.
//
//	  A secondary benefit of defending all users' freedom is that
//	improvements made in alternate versions of the program, if they
//	receive widespread use, become available for other developers to
//	incorporate.  Many developers of free software are heartened and
//	encouraged by the resulting cooperation.  However, in the case of
//	software used on network servers, this result may fail to come about.
//	The GNU General Public License permits making a modified version and
//	letting the public access it on a server without ever releasing its
//	source code to the public.
//
//	  The GNU Affero General Public License is designed specifically to
//	ensure that, in such cases, the modified source code becomes available
//	to the community.  It requires the operator of a network server to
//	provide the source code of the modified version running there to the
//	users of that server.  Therefore, public use of a modified version, on
//	a publicly accessible server, gives the public access to the source
//	code of the modified version.
//
//	  An older license, called the Affero General Public License and
//	published by Affero, was designed to accomplish similar goals.  This is
//	a different license, not a version of the Affero GPL, but Affero has
//	released a new version of the Affero GPL which permits relicensing under
//	this license.
//
//	  The precise terms and conditions for copying, distribution and
//	modification follow.
//
//	                       TERMS AND CONDITIONS
//
//	  0. Definitions.
//
//	  "This License" refers to version 3 of the GNU Affero General Public License.
//
//	  "Copyright" also means copyright-like laws that apply to other kinds of
//	works, such as semiconductor masks.
//
//	  "The Program" refers to any copyrightable work licensed under this
//	License.  Each licensee is addressed as "you".  "Licensees" and
//	"recipients" may be individuals or organizations.
//
//	  To "modify" a work means to copy from or adapt all or part of the work
//	in a fashion requiring copyright permission, other than the making of an
//	exact copy.  The resulting work is called a "modified version" of the
//	earlier work or a work "based on" the earlier work.
//
//	  A "covered work" means either the unmodified Program or a work based
//	on the Program.
//
//	  To "propagate" a work means to do anything with it that, without
//	permission, would make you directly or secondarily liable for
//	infringement under applicable copyright law, except executing it on a
//	computer or modifying a private copy.  Propagation includes copying,
//	distribution (with or without modification), making available to the
//	public, and in some countries other activities as well.
//
//	  To "convey" a work means any kind of propagation that enables other
//	parties to make or receive copies.  Mere interaction with a user through
//	a computer network, with no transfer of a copy, is not conveying.
//
//	  An interactive user interface displays "Appropriate Legal Notices"
//	to the extent that it includes a convenient and prominently visible
//	feature that (1) displays an appropriate copyright notice, and (2)
//	tells the user that there is no warranty for the work (except to the
//	extent that warranties are provided), that licensees may convey the
//	work under this License, and how to view a copy of this License.  If
//	the interface presents a list of user commands or options, such as a
//	menu, a prominent item in the list meets this criterion.
//
//	  1. Source Code.
//
//	  The "source code" for a work means the preferred form of the work
//	for making modifications to it.  "Object code" means any non-source
//	form of a work.
//
//	  A "Standard Interface" means an interface that either is an official
//	standard defined by a recognized standards body, or, in the case of
//	interfaces specified for a particular programming language, one that
//	is widely used among developers working in that language.
//
//	  The "System Libraries" of an executable work include anything, other
//	than the work as a whole, that (a) is included in the normal form of
//	packaging a Major Component, but which is not part of that Major
//	Component, and (b) serves only to enable use of the work with that
//	Major Component, or to implement a Standard Interface for which an
//	implementation is available to the public in source code form.  A
//	"Major Component", in this context, means a major essential component
//	(kernel, window system, and so on) of the specific operating system
//	(if any) on which the executable work runs, or a compiler used to
//	produce the work, or an object code interpreter used to run it.
//
//	  The "Corresponding Source" for a work in object code form means all
//	the source code needed to generate, install, and (for an executable
//	work) run the object code and to modify the work, including scripts to
//	control those activities.  However, it does not include the work's
//	System Libraries, or general-purpose tools or generally available free
//	programs which are used unmodified in performing those activities but
//	which are not part of the work.  For example, Corresponding Source
//	includes interface definition files associated with source files for
//	the work, and the source code for shared libraries and dynamically
//	linked subprograms that the work is specifically designed to require,
//	such as by intimate data communication or control flow between those
//	subprograms and other parts of the work.
//
//	  The Corresponding Source need not include anything that users
//	can regenerate automatically from other parts of the Corresponding
//	Source.
//
//	  The Corresponding Source for a work in source code form is that
//	same work.
//
//	  2. Basic Permissions.
//
//	  All rights granted under this License are granted for the term of
//	copyright on the Program, and are irrevocable provided the stated
//	conditions are met.  This License explicitly affirms your unlimited
//	permission to run the unmodified Program.  The output from running a
//	covered work is covered by this License only if the output, given its
//	content, constitutes a covered work.  This License acknowledges your
//	rights of fair use or other equivalent, as provided by copyright law.
//
//	  You may make, run and propagate covered works that you do not
//	convey, without conditions so long as your license otherwise remains
//	in force.  You may convey covered works to others for the sole purpose
//	of having them make modifications exclusively for you, or provide you
//	with facilities for running those works, provided that you comply with
//	the terms of this License in conveying all material for which you do
//	not control copyright.  Those thus making or running the covered works
//	for you must do so exclusively on your behalf, under your direction
//	and control, on terms that prohibit them from making any copies of
//	your copyrighted material outside their relationship with you.
//
//	  Conveying under any other circumstances is permitted solely under
//	the conditions stated below.  Sublicensing is not allowed; section 10
//	makes it unnecessary.
//
//	  3. Protecting Users' Legal Rights From Anti-Circumvention Law.
//
//	  No covered work shall be deemed part of an effective technological
//	measure under any applicable law fulfilling obligations under article
//	11 of the WIPO copyright treaty adopted on 20 December 1996, or
//	similar laws prohibiting or restricting circumvention of such
//	measures.
//
//	  When you convey a covered work, you waive any legal power to forbid
//	circumvention of technological measures to the extent such circumvention
//	is effected by exercising rights under this License with respect to
//	the covered work, and you disclaim any intention to limit operation or
//	modification of the work as a means of enforcing, against the work's
//	users, your or third parties' legal rights to forbid circumvention of
//	technological measures.
//
//	  4. Conveying Verbatim Copies.
//
//	  You may convey verbatim copies of the Program's source code as you
//	receive it, in any medium, provided that you conspicuously and
//	appropriately publish on each copy an appropriate copyright notice;
//	keep intact all notices stating that this License and any
//	non-permissive terms added in accord with section 7 apply to the code;
//	keep intact all notices of the absence of any warranty; and give all
//	recipients a copy of this License along with the Program.
//
//	  You may charge any price or no price for each copy that you convey,
//	and you may offer support or warranty protection for a fee.
//
//	  5. Conveying Modified Source Versions.
//
//	  You may convey a work based on the Program, or the modifications to
//	produce it from the Program, in the form of source code under the
//	terms of section 4, provided that you also meet all of these conditions:
//
//	    a) The work must carry prominent notices stating that you modified
//	    it, and giving a relevant date.
//
//	    b) The work must carry prominent notices stating that it is
//	    released under this License and any conditions added under section
//	    7.  This requirement modifies the requirement in section 4 to
//	    "keep intact all notices".
//
//	    c) You must license the entire work, as a whole, under this
//	    License to anyone who comes into possession of a copy.  This
//	    License will therefore apply, along with any applicable section 7
//	    additional terms, to the whole of the work, and all its parts,
//	    regardless of how they are packaged.  This License gives no
//	    permission to license the work in any other way, but it does not
//	    invalidate such permission if you have separately received it.
//
//	    d) If the work has interactive user interfaces, each must display
//	    Appropriate Legal Notices; however, if the Program has interactive
//	    interfaces that do not display Appropriate Legal Notices, your
//	    work need not make them do so.
//
//	  A compilation of a covered work with other separate and independent
//	works, which are not by their nature extensions of the covered work,
//	and which are not combined with it such as to form a larger program,
//	in or on a volume of a storage or distribution medium, is called an
//	"aggregate" if the compilation and its resulting copyright are not
//	used to limit the access or legal rights of the compilation's users
//	beyond what the individual works permit.  Inclusion of a covered work
//	in an aggregate does not cause this License to apply to the other
//	parts of the aggregate.
//
//	  6. Conveying Non-Source Forms.
//
//	  You may convey a covered work in object code form under the terms
//	of sections 4 and 5, provided that you also convey the
//	machine-readable Corresponding Source under the terms of this License,
//	in one of these ways:
//
//	    a) Convey the object code in, or embodied in, a physical product
//	    (including a physical distribution medium), accompanied by the
//	    Corresponding Source fixed on a durable physical medium
//	    customarily used for software interchange.
//
//	    b) Convey the object code in, or embodied in, a physical product
//	    (including a physical distribution medium), accompanied by a
//	    written offer, valid for at least three years and valid for as
//	    long as you offer spare parts or customer support for that product
//	    model, to give anyone who possesses the object code either (1) a
//	    copy of the Corresponding Source for all the software in the
//	    product that is covered by this License, on a durable physical
//	    medium customarily used for software interchange, for a price no
//	    more than your reasonable cost of physically performing this
//	    conveying of source, or (2) access to copy the
//	    Corresponding Source from a network server at no charge.
//
//	    c) Convey individual copies of the object code with a copy of the
//	    written offer to provide the Corresponding Source.  This
//	    alternative is allowed only occasionally and noncommercially, and
//	    only if you received the object code with such an offer, in accord
//	    with subsection 6b.
//
//	    d) Convey the object code by offering access from a designated
//	    place (gratis or for a charge), and offer equivalent access to the
//	    Corresponding Source in the same way through the same place at no
//	    further charge.  You need not require recipients to copy the
//	    Corresponding Source along with the object code.  If the place to
//	    copy the object code is a network server, the Corresponding Source
//	    may be on a different server (operated by you or a third party)
//	    that supports equivalent copying facilities, provided you maintain
//	    clear directions next to the object code saying where to find the
//	    Corresponding Source.  Regardless of what server hosts the
//	    Corresponding Source, you remain obligated to ensure that it is
//	    available for as long as needed to satisfy these requirements.
//
//	    e) Convey the object code using peer-to-peer transmission, provided
//	    you inform other peers where the object code and Corresponding
//	    Source of the work are being offered to the general public at no
//	    charge under subsection 6d.
//
//	  A separable portion of the object code, whose source code is excluded
//	from the Corresponding Source as a System Library, need not be
//	included in conveying the object code work.
//
//	  A "User Product" is either (1) a "consumer product", which means any
//	tangible personal property which is normally used for personal, family,
//	or household purposes, or (2) anything designed or sold for incorporation
//	into a dwelling.  In determining whether a product is a consumer product,
//	doubtful cases shall be resolved in favor of coverage.  For a particular
//	product received by a particular user, "normally used" refers to a
//	typical or common use of that class of product, regardless of the status
//	of the particular user or of the way in which the particular user
//	actually uses, or expects or is expected to use, the product.  A product
//	is a consumer product regardless of whether the product has substantial
//	commercial, industrial or non-consumer uses, unless such uses represent
//	the only significant mode of use of the product.
//
//	  "Installation Information" for a User Product means any methods,
//	procedures, authorization keys, or other information required to install
//	and execute modified versions of a covered work in that User Product from
//	a modified version of its Corresponding Source.  The information must
//	suffice to ensure that the continued functioning of the modified object
//	code is in no case prevented or interfered with solely because
//	modification has been made.
//
//	  If you convey an object code work under this section in, or with, or
//	specifically for use in, a User Product, and the conveying occurs as
//	part of a transaction in which the right of possession and use of the
//	User Product is transferred to the recipient in perpetuity or for a
//	fixed term (regardless of how the transaction is characterized), the
//	Corresponding Source conveyed under this section must be accompanied
//	by the Installation Information.  But this requirement does not apply
//	if neither you nor any third party retains the ability to install
//	modified object code on the User Product (for example, the work has
//	been installed in ROM).
//
//	  The requirement to provide Installation Information does not include a
//	requirement to continue to provide support service, warranty, or updates
//	for a work that has been modified or installed by the recipient, or for
//	the User Product in which it has been modified or installed.  Access to a
//	network may be denied when the modification itself materially and
//	adversely affects the operation of the network or violates the rules and
//	protocols for communication across the network.
//
//	  Corresponding Source conveyed, and Installation Information provided,
//	in accord with this section must be in a format that is publicly
//	documented (and with an implementation available to the public in
//	source code form), and must require no special password or key for
//	unpacking, reading or copying.
//
//	  7. Additional Terms.
//
//	  "Additional permissions" are terms that supplement the terms of this
//	License by making exceptions from one or more of its conditions.
//	Additional permissions that are applicable to the entire Program shall
//	be treated as though they were included in this License, to the extent
//	that they are valid under applicable law.  If additional permissions
//	apply only to part of the Program, that part may be used separately
//	under those permissions, but the entire Program remains governed by
//	this License without regard to the additional permissions.
//
//	  When you convey a copy of a covered work, you may at your option
//	remove any additional permissions from that copy, or from any part of
//	it.  (Additional permissions may be written to require their own
//	removal in certain cases when you modify the work.)  You may place
//	additional permissions on material, added by you to a covered work,
//	for which you have or can give appropriate copyright permission.
//
//	  Notwithstanding any other provision of this License, for material you
//	add to a covered work, you may (if authorized by the copyright holders of
//	that material) supplement the terms of this License with terms:
//
//	    a) Disclaiming warranty or limiting liability differently from the
//	    terms of sections 15 and 16 of this License; or
//
//	    b) Requiring preservation of specified reasonable legal notices or
//	    author attributions in that material or in the Appropriate Legal
//	    Notices displayed by works containing it; or
//
//	    c) Prohibiting misrepresentation of the origin of that material, or
//	    requiring that modified versions of such material be marked in
//	    reasonable ways as different from the original version; or
//
//	    d) Limiting the use for publicity purposes of names of licensors or
//	    authors of the material; or
//
//	    e) Declining to grant rights under trademark law for use of some
//	    trade names, trademarks, or service marks; or
//
//	    f) Requiring indemnification of licensors and authors of that
//	    material by anyone who conveys the material (or modified versions of
//	    it) with contractual assumptions of liability to the recipient, for
//	    any liability that these contractual assumptions directly impose on
//	    those licensors and authors.
//
//	  All other non-permissive additional terms are considered "further
//	restrictions" within the meaning of section 10.  If the Program as you
//	received it, or any part of it, contains a notice stating that it is
//	governed by this License along with a term that is a further
//	restriction, you may remove that term.  If a license document contains
//	a further restriction but permits relicensing or conveying under this
//	License, you may add to a covered work material governed by the terms
//	of that license document, provided that the further restriction does
//	not survive such relicensing or conveying.
//
//	  If you add terms to a covered work in accord with this section, you
//	must place, in the relevant source files, a statement of the
//	additional terms that apply to those files, or a notice indicating
//	where to find the applicable terms.
//
//	  Additional terms, permissive or non-permissive, may be stated in the
//	form of a separately written license, or stated as exceptions;
//	the above requirements apply either way.
//
//	  8. Termination.
//
//	  You may not propagate or modify a covered work except as expressly
//	provided under this License.  Any attempt otherwise to propagate or
//	modify it is void, and will automatically terminate your rights under
//	this License (including any patent licenses granted under the third
//	paragraph of section 11).
//
//	  However, if you cease all violation of this License, then your
//	license from a particular copyright holder is reinstated (a)
//	provisionally, unless and until the copyright holder explicitly and
//	finally terminates your license, and (b) permanently, if the copyright
//	holder fails to notify you of the violation by some reasonable means
//	prior to 60 days after the cessation.
//
//	  Moreover, your license from a particular copyright holder is
//	reinstated permanently if the copyright holder notifies you of the
//	violation by some reasonable means, this is the first time you have
//	received notice of violation of this License (for any work) from that
//	copyright holder, and you cure the violation prior to 30 days after
//	your receipt of the notice.
//
//	  Termination of your rights under this section does not terminate the
//	licenses of parties who have received copies or rights from you under
//	this License.  If your rights have been terminated and not permanently
//	reinstated, you do not qualify to receive new licenses for the same
//	material under section 10.
//
//	  9. Acceptance Not Required for Having Copies.
//
//	  You are not required to accept this License in order to receive or
//	run a copy of the Program.  Ancillary propagation of a covered work
//	occurring solely as a consequence of using peer-to-peer transmission
//	to receive a copy likewise does not require acceptance.  However,
//	nothing other than this License grants you permission to propagate or
//	modify any covered work.  These actions infringe copyright if you do
//	not accept this License.  Therefore, by modifying or propagating a
//	covered work, you indicate your acceptance of this License to do so.
//
//	  10. Automatic Licensing of Downstream Recipients.
//
//	  Each time you convey a covered work, the recipient automatically
//	receives a license from the original licensors, to run, modify and
//	propagate that work, subject to this License.  You are not responsible
//	for enforcing compliance by third parties with this License.
//
//	  An "entity transaction" is a transaction transferring control of an
//	organization, or substantially all assets of one, or subdividing an
//	organization, or merging organizations.  If propagation of a covered
//	work results from an entity transaction, each party to that
//	transaction who receives a copy of the work also receives whatever
//	licenses to the work the party's predecessor in interest had or could
//	give under the previous paragraph, plus a right to possession of the
//	Corresponding Source of the work from the predecessor in interest, if
//	the predecessor has it or can get it with reasonable efforts.
//
//	  You may not impose any further restrictions on the exercise of the
//	rights granted or affirmed under this License.  For example, you may
//	not impose a license fee, royalty, or other charge for exercise of
//	rights granted under this License, and you may not initiate litigation
//	(including a cross-claim or counterclaim in a lawsuit) alleging that
//	any patent claim is infringed by making, using, selling, offering for
//	sale, or importing the Program or any portion of it.
//
//	  11. Patents.
//
//	  A "contributor" is a copyright holder who authorizes use under this
//	License of the Program or a work on which the Program is based.  The
//	work thus licensed is called the contributor's "contributor version".
//
//	  A contributor's "essential patent claims" are all patent claims
//	owned or controlled by the contributor, whether already acquired or
//	hereafter acquired, that would be infringed by some manner, permitted
//	by this License, of making, using, or selling its contributor version,
//	but do not include claims that would be infringed only as a
//	consequence of further modification of the contributor version.  For
//	purposes of this definition, "control" includes the right to grant
//	patent sublicenses in a manner consistent with the requirements of
//	this License.
//
//	  Each contributor grants you a non-exclusive, worldwide, royalty-free
//	patent license under the contributor's essential patent claims, to
//	make, use, sell, offer for sale, import and otherwise run, modify and
//	propagate the contents of its contributor version.
//
//	  In the following three paragraphs, a "patent license" is any express
//	agreement or commitment, however denominated, not to enforce a patent
//	(such as an express permission to practice a patent or covenant not to
//	sue for patent infringement).  To "grant" such a patent license to a
//	party means to make such an agreement or commitment not to enforce a
//	patent against the party.
//
//	  If you convey a covered work, knowingly relying on a patent license,
//	and the Corresponding Source of the work is not available for anyone
//	to copy, free of charge and under the terms of this License, through a
//	publicly available network server or other readily accessible means,
//	then you must either (1) cause the Corresponding Source to be so
//	available, or (2) arrange to deprive yourself of the benefit of the
//	patent license for this particular work, or (3) arrange, in a manner
//	consistent with the requirements of this License, to extend the patent
//	license to downstream recipients.  "Knowingly relying" means you have
//	actual knowledge that, but for the patent license, your conveying the
//	covered work in a country, or your recipient's use of the covered work
//	in a country, would infringe one or more identifiable patents in that
//	country that you have reason to believe are valid.
//
//	  If, pursuant to or in connection with a single transaction or
//	arrangement, you convey, or propagate by procuring conveyance of, a
//	covered work, and grant a patent license to some of the parties
//	receiving the covered work authorizing them to use, propagate, modify
//	or convey a specific copy of the covered work, then the patent license
//	you grant is automatically extended to all recipients of the covered
//	work and works based on it.
//
//	  A patent license is "discriminatory" if it does not include within
//	the scope of its coverage, prohibits the exercise of, or is
//	conditioned on the non-exercise of one or more of the rights that are
//	specifically granted under this License.  You may not convey a covered
//	work if you are a party to an arrangement with a third party that is
//	in the business of distributing software, under which you make payment
//	to the third party based on the extent of your activity of conveying
//	the work, and under which the third party grants, to any of the
//	parties who would receive the covered work from you, a discriminatory
//	patent license (a) in connection with copies of the covered work
//	conveyed by you (or copies made from those copies), or (b) primarily
//	for and in connection with specific products or compilations that
//	contain the covered work, unless you entered into that arrangement,
//	or that patent license was granted, prior to 28 March 2007.
//
//	  Nothing in this License shall be construed as excluding or limiting
//	any implied license or other defenses to infringement that may
//	otherwise be available to you under applicable patent law.
//
//	  12. No Surrender of Others' Freedom.
//
//	  If conditions are imposed on you (whether by court order, agreement or
//	otherwise) that contradict the conditions of this License, they do not
//	excuse you from the conditions of this License.  If you cannot convey a
//	covered work so as to satisfy simultaneously your obligations under this
//	License and any other pertinent obligations, then as a consequence you may
//	not convey it at all.  For example, if you agree to terms that obligate you
//	to collect a royalty for further conveying from those to whom you convey
//	the Program, the only way you could satisfy both those terms and this
//	License would be to refrain entirely from conveying the Program.
//
//	  13. Remote Network Interaction; Use with the GNU General Public License.
//
//	  Notwithstanding any other provision of this License, if you modify the
//	Program, your modified version must prominently offer all users
//	interacting with it remotely through a computer network (if your version
//	supports such interaction) an opportunity to receive the Corresponding
//	Source of your version by providing access to the Corresponding Source
//	from a network server at no charge, through some standard or customary
//	means of facilitating copying of software.  This Corresponding Source
//	shall include the Corresponding Source for any work covered by version 3
//	of the GNU General Public License that is incorporated pursuant to the
//	following paragraph.
//
//	  Notwithstanding any other provision of this License, you have
//	permission to link or combine any covered work with a work licensed
//	under version 3 of the GNU General Public License into a single
//	combined work, and to convey the resulting work.  The terms of this
//	License will continue to apply to the part which is the covered work,
//	but the work with which it is combined will remain governed by version
//	3 of the GNU General Public License.
//
//	  14. Revised Versions of this License.
//
//	  The Free Software Foundation may publish revised and/or new versions of
//	the GNU Affero General Public License from time to time.  Such new versions
//	will be similar in spirit to the present version, but may differ in detail to
//	address new problems or concerns.
//
//	  Each version is given a distinguishing version number.  If the
//	Program specifies that a certain numbered version of the GNU Affero General
//	Public License "or any later version" applies to it, you have the
//	option of following the terms and conditions either of that numbered
//	version or of any later version published by the Free Software
//	Foundation.  If the Program does not specify a version number of the
//	GNU Affero General Public License, you may choose any version ever published
//	by the Free Software Foundation.
//
//	  If the Program specifies that a proxy can decide which future
//	versions of the GNU Affero General Public License can be used, that proxy's
//	public statement of acceptance of a version permanently authorizes you
//	to choose that version for the Program.
//
//	  Later license versions may give you additional or different
//	permissions.  However, no additional obligations are imposed on any
//	author or copyright holder as a result of your choosing to follow a
//	later version.
//
//	  15. Disclaimer of Warranty.
//
//	  THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY
//	APPLICABLE LAW.  EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT
//	HOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY
//	OF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
//	THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
//	PURPOSE.  THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM
//	IS WITH YOU.  SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF
//	ALL NECESSARY SERVICING, REPAIR OR CORRECTION.
//
//	  16. Limitation of Liability.
//
//	  IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING
//	WILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MODIFIES AND/OR CONVEYS
//	THE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES, INCLUDING ANY
//	GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF THE
//	USE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED TO LOSS OF
//	DATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD
//	PARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER PROGRAMS),
//	EVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF
//	SUCH DAMAGES.
//
//	  17. Interpretation of Sections 15 and 16.
//
//	  If the disclaimer of warranty and limitation of liability provided
//	above cannot be given local legal effect according to their terms,
//	reviewing courts shall apply local law that most closely approximates
//	an absolute waiver of all civil liability in connection with the
//	Program, unless a warranty or assumption of liability accompanies a
//	copy of the Program in return for a fee.
//
//	                     END OF TERMS AND CONDITIONS
//
//	            How to Apply These Terms to Your New Programs
//
//	  If you develop a new program, and you want it to be of the greatest
//	possible use to the public, the best way to achieve this is to make it
//	free software which everyone can redistribute and change under these terms.
//
//	  To do so, attach the following notices to the program.  It is safest
//	to attach them to the start of each source file to most effectively
//	state the exclusion of warranty; and each file should have at least
//	the "copyright" line and a pointer to where the full notice is found.
//
//	    <one line to give the program's name and a brief idea of what it does.>
//	    Copyright (C) <year>  <name of author>
//
//	    This program is free software: you can redistribute it and/or modify
//	    it under the terms of the GNU Affero General Public License as published by
//	    the Free Software Foundation, either version 3 of the License, or
//	    (at your option) any later version.
//
//	    This program is distributed in the hope that it will be useful,
//	    but WITHOUT ANY WARRANTY; without even the implied warranty of
//	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	    GNU Affero General Public License for more details.
//
//	    You should have received a copy of the GNU Affero General Public License
//	    along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
//	Also add information on how to contact you by electronic and paper mail.
//
//	  If your software can interact with users remotely through a computer
//	network, you should also make sure that it provides a way for users to
//	get its source.  For example, if your program is a web application, its
//	interface could display a "Source" link that leads users to an archive
//	of the code.  There are many ways you could offer source, and different
//	solutions will be better for different programs; see section 13 for the
//	specific requirements.
//
//	  You should also get your employer (if you work as a programmer) or school,
//	if any, to sign a "copyright disclaimer" for the program, if necessary.
//	For more information on this, and how to apply and follow the GNU AGPL, see
//	<https://www.gnu.org/licenses/>.
//
// # Generation
//
// This variable and the associated documentation have been automatically generated using the 'gogenlicense' tool.
// It was last updated at 15-01-2023 11:33:49.
var LegalNotices string

func init() {
	LegalNotices = "The following go packages are imported:\n- Go Standard Library (BSD-3-Clause; see https://golang.org/LICENSE)\n- gorm.io/gorm (MIT; see https://github.com/go-gorm/gorm/blob/v1.24.2/License)\n- gorm.io/driver/mysql (MIT; see https://github.com/go-gorm/mysql/blob/v1.4.4/License)\n- gopkg.in/yaml.v2 (Apache-2.0; see https://github.com/go-yaml/yaml/blob/v2.3.0/LICENSE)\n- golang.org/x/tools/go/analysis (BSD-3-Clause; see https://cs.opensource.google/go/x/tools/+/v0.4.0:LICENSE)\n- golang.org/x/term (BSD-3-Clause; see https://cs.opensource.google/go/x/term/+/v0.3.0:LICENSE)\n- golang.org/x/sys/unix (BSD-3-Clause; see https://cs.opensource.google/go/x/sys/+/v0.3.0:LICENSE)\n- golang.org/x/sync/errgroup (BSD-3-Clause; see https://cs.opensource.google/go/x/sync/+/v0.1.0:LICENSE)\n- golang.org/x/exp (BSD-3-Clause; see https://cs.opensource.google/go/x/exp/+/8879d019:LICENSE)\n- golang.org/x/crypto (BSD-3-Clause; see https://cs.opensource.google/go/x/crypto/+/v0.3.0:LICENSE)\n- github.com/yuin/goldmark-meta (MIT; see https://github.com/yuin/goldmark-meta/blob/v1.1.0/LICENSE)\n- github.com/yuin/goldmark (MIT; see https://github.com/yuin/goldmark/blob/v1.4.13/LICENSE)\n- github.com/tkw1536/goprogram (MIT; see https://github.com/tkw1536/goprogram/blob/v0.2.4/LICENSE)\n- github.com/tdewolff/parse (MIT; see https://github.com/tdewolff/parse/blob/v2.3.4/LICENSE.md)\n- github.com/tdewolff/minify (MIT; see https://github.com/tdewolff/minify/blob/v2.3.6/LICENSE.md)\n- github.com/rs/zerolog (MIT; see https://github.com/rs/zerolog/blob/v1.28.0/LICENSE)\n- github.com/pquerna/otp (Apache-2.0; see https://github.com/pquerna/otp/blob/v1.4.0/LICENSE)\n- github.com/pkg/errors (BSD-2-Clause; see https://github.com/pkg/errors/blob/v0.9.1/LICENSE)\n- github.com/mattn/go-isatty (MIT; see https://github.com/mattn/go-isatty/blob/v0.0.16/LICENSE)\n- github.com/mattn/go-colorable (MIT; see https://github.com/mattn/go-colorable/blob/v0.1.13/LICENSE)\n- github.com/julienschmidt/httprouter (BSD-3-Clause; see https://github.com/julienschmidt/httprouter/blob/v1.3.0/LICENSE)\n- github.com/jinzhu/now (MIT; see https://github.com/jinzhu/now/blob/v1.1.5/License)\n- github.com/jinzhu/inflection (MIT; see https://github.com/jinzhu/inflection/blob/v1.0.0/LICENSE)\n- github.com/jessevdk/go-flags (BSD-3-Clause; see https://github.com/jessevdk/go-flags/blob/v1.5.0/LICENSE)\n- github.com/gosuri/uilive (MIT; see https://github.com/gosuri/uilive/blob/v0.0.4/LICENSE)\n- github.com/gorilla/websocket (BSD-2-Clause; see https://github.com/gorilla/websocket/blob/v1.5.0/LICENSE)\n- github.com/gorilla/sessions (BSD-3-Clause; see https://github.com/gorilla/sessions/blob/v1.2.1/LICENSE)\n- github.com/gorilla/securecookie (BSD-3-Clause; see https://github.com/gorilla/securecookie/blob/v1.1.1/LICENSE)\n- github.com/gorilla/csrf (BSD-3-Clause; see https://github.com/gorilla/csrf/blob/v1.7.1/LICENSE)\n- github.com/go-sql-driver/mysql (MPL-2.0; see https://github.com/go-sql-driver/mysql/blob/v1.6.0/LICENSE)\n- github.com/gliderlabs/ssh (BSD-3-Clause; see https://github.com/gliderlabs/ssh/blob/v0.3.5/LICENSE)\n- github.com/feiin/sqlstring (MIT; see https://github.com/feiin/sqlstring/blob/v0.3.0/LICENSE)\n- github.com/boombuler/barcode (MIT; see https://github.com/boombuler/barcode/blob/6c824513bacc/LICENSE)\n- github.com/anmitsu/go-shlex (MIT; see https://github.com/anmitsu/go-shlex/blob/38f4b401e2be/LICENSE)\n- github.com/alessio/shellescape (MIT; see https://github.com/alessio/shellescape/blob/v1.4.1/LICENSE)\n- github.com/Showmax/go-fqdn (Apache-2.0; see https://github.com/Showmax/go-fqdn/blob/v1.0.0/LICENSE)\n- github.com/FAU-CDI/wdresolve (AGPL-3.0; see https://github.com/FAU-CDI/wdresolve/blob/c9c6779d7c41/LICENSE)\n\n================================================================================\n\n\n================================================================================\nGo Standard Library\nLicensed under the Terms of the BSD-3-Clause License, see also https://golang.org/LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n\t* Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n\t* Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n\t* Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule gorm.io/gorm\nLicensed under the Terms of the MIT License, see also https://github.com/go-gorm/gorm/blob/v1.24.2/License. \n\nThe MIT License (MIT)\n\nCopyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in\nall copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN\nTHE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule gorm.io/driver/mysql\nLicensed under the Terms of the MIT License, see also https://github.com/go-gorm/mysql/blob/v1.4.4/License. \n\nThe MIT License (MIT)\n\nCopyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in\nall copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN\nTHE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule gopkg.in/yaml.v2\nLicensed under the Terms of the Apache-2.0 License, see also https://github.com/go-yaml/yaml/blob/v2.3.0/LICENSE. \n\n                                 Apache License\n                           Version 2.0, January 2004\n                        http://www.apache.org/licenses/\n\n   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION\n\n   1. Definitions.\n\n      \"License\" shall mean the terms and conditions for use, reproduction,\n      and distribution as defined by Sections 1 through 9 of this document.\n\n      \"Licensor\" shall mean the copyright owner or entity authorized by\n      the copyright owner that is granting the License.\n\n      \"Legal Entity\" shall mean the union of the acting entity and all\n      other entities that control, are controlled by, or are under common\n      control with that entity. For the purposes of this definition,\n      \"control\" means (i) the power, direct or indirect, to cause the\n      direction or management of such entity, whether by contract or\n      otherwise, or (ii) ownership of fifty percent (50%) or more of the\n      outstanding shares, or (iii) beneficial ownership of such entity.\n\n      \"You\" (or \"Your\") shall mean an individual or Legal Entity\n      exercising permissions granted by this License.\n\n      \"Source\" form shall mean the preferred form for making modifications,\n      including but not limited to software source code, documentation\n      source, and configuration files.\n\n      \"Object\" form shall mean any form resulting from mechanical\n      transformation or translation of a Source form, including but\n      not limited to compiled object code, generated documentation,\n      and conversions to other media types.\n\n      \"Work\" shall mean the work of authorship, whether in Source or\n      Object form, made available under the License, as indicated by a\n      copyright notice that is included in or attached to the work\n      (an example is provided in the Appendix below).\n\n      \"Derivative Works\" shall mean any work, whether in Source or Object\n      form, that is based on (or derived from) the Work and for which the\n      editorial revisions, annotations, elaborations, or other modifications\n      represent, as a whole, an original work of authorship. For the purposes\n      of this License, Derivative Works shall not include works that remain\n      separable from, or merely link (or bind by name) to the interfaces of,\n      the Work and Derivative Works thereof.\n\n      \"Contribution\" shall mean any work of authorship, including\n      the original version of the Work and any modifications or additions\n      to that Work or Derivative Works thereof, that is intentionally\n      submitted to Licensor for inclusion in the Work by the copyright owner\n      or by an individual or Legal Entity authorized to submit on behalf of\n      the copyright owner. For the purposes of this definition, \"submitted\"\n      means any form of electronic, verbal, or written communication sent\n      to the Licensor or its representatives, including but not limited to\n      communication on electronic mailing lists, source code control systems,\n      and issue tracking systems that are managed by, or on behalf of, the\n      Licensor for the purpose of discussing and improving the Work, but\n      excluding communication that is conspicuously marked or otherwise\n      designated in writing by the copyright owner as \"Not a Contribution.\"\n\n      \"Contributor\" shall mean Licensor and any individual or Legal Entity\n      on behalf of whom a Contribution has been received by Licensor and\n      subsequently incorporated within the Work.\n\n   2. Grant of Copyright License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      copyright license to reproduce, prepare Derivative Works of,\n      publicly display, publicly perform, sublicense, and distribute the\n      Work and such Derivative Works in Source or Object form.\n\n   3. Grant of Patent License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      (except as stated in this section) patent license to make, have made,\n      use, offer to sell, sell, import, and otherwise transfer the Work,\n      where such license applies only to those patent claims licensable\n      by such Contributor that are necessarily infringed by their\n      Contribution(s) alone or by combination of their Contribution(s)\n      with the Work to which such Contribution(s) was submitted. If You\n      institute patent litigation against any entity (including a\n      cross-claim or counterclaim in a lawsuit) alleging that the Work\n      or a Contribution incorporated within the Work constitutes direct\n      or contributory patent infringement, then any patent licenses\n      granted to You under this License for that Work shall terminate\n      as of the date such litigation is filed.\n\n   4. Redistribution. You may reproduce and distribute copies of the\n      Work or Derivative Works thereof in any medium, with or without\n      modifications, and in Source or Object form, provided that You\n      meet the following conditions:\n\n      (a) You must give any other recipients of the Work or\n          Derivative Works a copy of this License; and\n\n      (b) You must cause any modified files to carry prominent notices\n          stating that You changed the files; and\n\n      (c) You must retain, in the Source form of any Derivative Works\n          that You distribute, all copyright, patent, trademark, and\n          attribution notices from the Source form of the Work,\n          excluding those notices that do not pertain to any part of\n          the Derivative Works; and\n\n      (d) If the Work includes a \"NOTICE\" text file as part of its\n          distribution, then any Derivative Works that You distribute must\n          include a readable copy of the attribution notices contained\n          within such NOTICE file, excluding those notices that do not\n          pertain to any part of the Derivative Works, in at least one\n          of the following places: within a NOTICE text file distributed\n          as part of the Derivative Works; within the Source form or\n          documentation, if provided along with the Derivative Works; or,\n          within a display generated by the Derivative Works, if and\n          wherever such third-party notices normally appear. The contents\n          of the NOTICE file are for informational purposes only and\n          do not modify the License. You may add Your own attribution\n          notices within Derivative Works that You distribute, alongside\n          or as an addendum to the NOTICE text from the Work, provided\n          that such additional attribution notices cannot be construed\n          as modifying the License.\n\n      You may add Your own copyright statement to Your modifications and\n      may provide additional or different license terms and conditions\n      for use, reproduction, or distribution of Your modifications, or\n      for any such Derivative Works as a whole, provided Your use,\n      reproduction, and distribution of the Work otherwise complies with\n      the conditions stated in this License.\n\n   5. Submission of Contributions. Unless You explicitly state otherwise,\n      any Contribution intentionally submitted for inclusion in the Work\n      by You to the Licensor shall be under the terms and conditions of\n      this License, without any additional terms or conditions.\n      Notwithstanding the above, nothing herein shall supersede or modify\n      the terms of any separate license agreement you may have executed\n      with Licensor regarding such Contributions.\n\n   6. Trademarks. This License does not grant permission to use the trade\n      names, trademarks, service marks, or product names of the Licensor,\n      except as required for reasonable and customary use in describing the\n      origin of the Work and reproducing the content of the NOTICE file.\n\n   7. Disclaimer of Warranty. Unless required by applicable law or\n      agreed to in writing, Licensor provides the Work (and each\n      Contributor provides its Contributions) on an \"AS IS\" BASIS,\n      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or\n      implied, including, without limitation, any warranties or conditions\n      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A\n      PARTICULAR PURPOSE. You are solely responsible for determining the\n      appropriateness of using or redistributing the Work and assume any\n      risks associated with Your exercise of permissions under this License.\n\n   8. Limitation of Liability. In no event and under no legal theory,\n      whether in tort (including negligence), contract, or otherwise,\n      unless required by applicable law (such as deliberate and grossly\n      negligent acts) or agreed to in writing, shall any Contributor be\n      liable to You for damages, including any direct, indirect, special,\n      incidental, or consequential damages of any character arising as a\n      result of this License or out of the use or inability to use the\n      Work (including but not limited to damages for loss of goodwill,\n      work stoppage, computer failure or malfunction, or any and all\n      other commercial damages or losses), even if such Contributor\n      has been advised of the possibility of such damages.\n\n   9. Accepting Warranty or Additional Liability. While redistributing\n      the Work or Derivative Works thereof, You may choose to offer,\n      and charge a fee for, acceptance of support, warranty, indemnity,\n      or other liability obligations and/or rights consistent with this\n      License. However, in accepting such obligations, You may act only\n      on Your own behalf and on Your sole responsibility, not on behalf\n      of any other Contributor, and only if You agree to indemnify,\n      defend, and hold each Contributor harmless for any liability\n      incurred by, or claims asserted against, such Contributor by reason\n      of your accepting any such warranty or additional liability.\n\n   END OF TERMS AND CONDITIONS\n\n   APPENDIX: How to apply the Apache License to your work.\n\n      To apply the Apache License to your work, attach the following\n      boilerplate notice, with the fields enclosed by brackets \"{}\"\n      replaced with your own identifying information. (Don't include\n      the brackets!)  The text should be enclosed in the appropriate\n      comment syntax for the file format. We also recommend that a\n      file or class name and description of purpose be included on the\n      same \"printed page\" as the copyright notice for easier\n      identification within third-party archives.\n\n   Copyright {yyyy} {name of copyright owner}\n\n   Licensed under the Apache License, Version 2.0 (the \"License\");\n   you may not use this file except in compliance with the License.\n   You may obtain a copy of the License at\n\n       http://www.apache.org/licenses/LICENSE-2.0\n\n   Unless required by applicable law or agreed to in writing, software\n   distributed under the License is distributed on an \"AS IS\" BASIS,\n   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n   See the License for the specific language governing permissions and\n   limitations under the License.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/tools/go/analysis\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/tools/+/v0.4.0:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/term\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/term/+/v0.3.0:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/sys/unix\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/sys/+/v0.3.0:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/sync/errgroup\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/sync/+/v0.1.0:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/exp\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/exp/+/8879d019:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule golang.org/x/crypto\nLicensed under the Terms of the BSD-3-Clause License, see also https://cs.opensource.google/go/x/crypto/+/v0.3.0:LICENSE. \n\nCopyright (c) 2009 The Go Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/yuin/goldmark-meta\nLicensed under the Terms of the MIT License, see also https://github.com/yuin/goldmark-meta/blob/v1.1.0/LICENSE. \n\nMIT License\n\nCopyright (c) 2019 Yusuke Inuzuka\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/yuin/goldmark\nLicensed under the Terms of the MIT License, see also https://github.com/yuin/goldmark/blob/v1.4.13/LICENSE. \n\nMIT License\n\nCopyright (c) 2019 Yusuke Inuzuka\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/tkw1536/goprogram\nLicensed under the Terms of the MIT License, see also https://github.com/tkw1536/goprogram/blob/v0.2.4/LICENSE. \n\nMIT License\n\nCopyright (c) 2022 Tom Wiesing\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/tdewolff/parse\nLicensed under the Terms of the MIT License, see also https://github.com/tdewolff/parse/blob/v2.3.4/LICENSE.md. \n\nCopyright (c) 2015 Taco de Wolff\n\n Permission is hereby granted, free of charge, to any person\n obtaining a copy of this software and associated documentation\n files (the \"Software\"), to deal in the Software without\n restriction, including without limitation the rights to use,\n copy, modify, merge, publish, distribute, sublicense, and/or sell\n copies of the Software, and to permit persons to whom the\n Software is furnished to do so, subject to the following\n conditions:\n\n The above copyright notice and this permission notice shall be\n included in all copies or substantial portions of the Software.\n\n THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND,\n EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES\n OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND\n NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT\n HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,\n WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING\n FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR\n OTHER DEALINGS IN THE SOFTWARE.\n================================================================================\n\n================================================================================\nModule github.com/tdewolff/minify\nLicensed under the Terms of the MIT License, see also https://github.com/tdewolff/minify/blob/v2.3.6/LICENSE.md. \n\nCopyright (c) 2015 Taco de Wolff\n\n Permission is hereby granted, free of charge, to any person\n obtaining a copy of this software and associated documentation\n files (the \"Software\"), to deal in the Software without\n restriction, including without limitation the rights to use,\n copy, modify, merge, publish, distribute, sublicense, and/or sell\n copies of the Software, and to permit persons to whom the\n Software is furnished to do so, subject to the following\n conditions:\n\n The above copyright notice and this permission notice shall be\n included in all copies or substantial portions of the Software.\n\n THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND,\n EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES\n OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND\n NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT\n HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,\n WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING\n FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR\n OTHER DEALINGS IN THE SOFTWARE.\n================================================================================\n\n================================================================================\nModule github.com/rs/zerolog\nLicensed under the Terms of the MIT License, see also https://github.com/rs/zerolog/blob/v1.28.0/LICENSE. \n\nMIT License\n\nCopyright (c) 2017 Olivier Poitrey\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/pquerna/otp\nLicensed under the Terms of the Apache-2.0 License, see also https://github.com/pquerna/otp/blob/v1.4.0/LICENSE. \n\n\n                                 Apache License\n                           Version 2.0, January 2004\n                        http://www.apache.org/licenses/\n\n   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION\n\n   1. Definitions.\n\n      \"License\" shall mean the terms and conditions for use, reproduction,\n      and distribution as defined by Sections 1 through 9 of this document.\n\n      \"Licensor\" shall mean the copyright owner or entity authorized by\n      the copyright owner that is granting the License.\n\n      \"Legal Entity\" shall mean the union of the acting entity and all\n      other entities that control, are controlled by, or are under common\n      control with that entity. For the purposes of this definition,\n      \"control\" means (i) the power, direct or indirect, to cause the\n      direction or management of such entity, whether by contract or\n      otherwise, or (ii) ownership of fifty percent (50%) or more of the\n      outstanding shares, or (iii) beneficial ownership of such entity.\n\n      \"You\" (or \"Your\") shall mean an individual or Legal Entity\n      exercising permissions granted by this License.\n\n      \"Source\" form shall mean the preferred form for making modifications,\n      including but not limited to software source code, documentation\n      source, and configuration files.\n\n      \"Object\" form shall mean any form resulting from mechanical\n      transformation or translation of a Source form, including but\n      not limited to compiled object code, generated documentation,\n      and conversions to other media types.\n\n      \"Work\" shall mean the work of authorship, whether in Source or\n      Object form, made available under the License, as indicated by a\n      copyright notice that is included in or attached to the work\n      (an example is provided in the Appendix below).\n\n      \"Derivative Works\" shall mean any work, whether in Source or Object\n      form, that is based on (or derived from) the Work and for which the\n      editorial revisions, annotations, elaborations, or other modifications\n      represent, as a whole, an original work of authorship. For the purposes\n      of this License, Derivative Works shall not include works that remain\n      separable from, or merely link (or bind by name) to the interfaces of,\n      the Work and Derivative Works thereof.\n\n      \"Contribution\" shall mean any work of authorship, including\n      the original version of the Work and any modifications or additions\n      to that Work or Derivative Works thereof, that is intentionally\n      submitted to Licensor for inclusion in the Work by the copyright owner\n      or by an individual or Legal Entity authorized to submit on behalf of\n      the copyright owner. For the purposes of this definition, \"submitted\"\n      means any form of electronic, verbal, or written communication sent\n      to the Licensor or its representatives, including but not limited to\n      communication on electronic mailing lists, source code control systems,\n      and issue tracking systems that are managed by, or on behalf of, the\n      Licensor for the purpose of discussing and improving the Work, but\n      excluding communication that is conspicuously marked or otherwise\n      designated in writing by the copyright owner as \"Not a Contribution.\"\n\n      \"Contributor\" shall mean Licensor and any individual or Legal Entity\n      on behalf of whom a Contribution has been received by Licensor and\n      subsequently incorporated within the Work.\n\n   2. Grant of Copyright License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      copyright license to reproduce, prepare Derivative Works of,\n      publicly display, publicly perform, sublicense, and distribute the\n      Work and such Derivative Works in Source or Object form.\n\n   3. Grant of Patent License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      (except as stated in this section) patent license to make, have made,\n      use, offer to sell, sell, import, and otherwise transfer the Work,\n      where such license applies only to those patent claims licensable\n      by such Contributor that are necessarily infringed by their\n      Contribution(s) alone or by combination of their Contribution(s)\n      with the Work to which such Contribution(s) was submitted. If You\n      institute patent litigation against any entity (including a\n      cross-claim or counterclaim in a lawsuit) alleging that the Work\n      or a Contribution incorporated within the Work constitutes direct\n      or contributory patent infringement, then any patent licenses\n      granted to You under this License for that Work shall terminate\n      as of the date such litigation is filed.\n\n   4. Redistribution. You may reproduce and distribute copies of the\n      Work or Derivative Works thereof in any medium, with or without\n      modifications, and in Source or Object form, provided that You\n      meet the following conditions:\n\n      (a) You must give any other recipients of the Work or\n          Derivative Works a copy of this License; and\n\n      (b) You must cause any modified files to carry prominent notices\n          stating that You changed the files; and\n\n      (c) You must retain, in the Source form of any Derivative Works\n          that You distribute, all copyright, patent, trademark, and\n          attribution notices from the Source form of the Work,\n          excluding those notices that do not pertain to any part of\n          the Derivative Works; and\n\n      (d) If the Work includes a \"NOTICE\" text file as part of its\n          distribution, then any Derivative Works that You distribute must\n          include a readable copy of the attribution notices contained\n          within such NOTICE file, excluding those notices that do not\n          pertain to any part of the Derivative Works, in at least one\n          of the following places: within a NOTICE text file distributed\n          as part of the Derivative Works; within the Source form or\n          documentation, if provided along with the Derivative Works; or,\n          within a display generated by the Derivative Works, if and\n          wherever such third-party notices normally appear. The contents\n          of the NOTICE file are for informational purposes only and\n          do not modify the License. You may add Your own attribution\n          notices within Derivative Works that You distribute, alongside\n          or as an addendum to the NOTICE text from the Work, provided\n          that such additional attribution notices cannot be construed\n          as modifying the License.\n\n      You may add Your own copyright statement to Your modifications and\n      may provide additional or different license terms and conditions\n      for use, reproduction, or distribution of Your modifications, or\n      for any such Derivative Works as a whole, provided Your use,\n      reproduction, and distribution of the Work otherwise complies with\n      the conditions stated in this License.\n\n   5. Submission of Contributions. Unless You explicitly state otherwise,\n      any Contribution intentionally submitted for inclusion in the Work\n      by You to the Licensor shall be under the terms and conditions of\n      this License, without any additional terms or conditions.\n      Notwithstanding the above, nothing herein shall supersede or modify\n      the terms of any separate license agreement you may have executed\n      with Licensor regarding such Contributions.\n\n   6. Trademarks. This License does not grant permission to use the trade\n      names, trademarks, service marks, or product names of the Licensor,\n      except as required for reasonable and customary use in describing the\n      origin of the Work and reproducing the content of the NOTICE file.\n\n   7. Disclaimer of Warranty. Unless required by applicable law or\n      agreed to in writing, Licensor provides the Work (and each\n      Contributor provides its Contributions) on an \"AS IS\" BASIS,\n      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or\n      implied, including, without limitation, any warranties or conditions\n      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A\n      PARTICULAR PURPOSE. You are solely responsible for determining the\n      appropriateness of using or redistributing the Work and assume any\n      risks associated with Your exercise of permissions under this License.\n\n   8. Limitation of Liability. In no event and under no legal theory,\n      whether in tort (including negligence), contract, or otherwise,\n      unless required by applicable law (such as deliberate and grossly\n      negligent acts) or agreed to in writing, shall any Contributor be\n      liable to You for damages, including any direct, indirect, special,\n      incidental, or consequential damages of any character arising as a\n      result of this License or out of the use or inability to use the\n      Work (including but not limited to damages for loss of goodwill,\n      work stoppage, computer failure or malfunction, or any and all\n      other commercial damages or losses), even if such Contributor\n      has been advised of the possibility of such damages.\n\n   9. Accepting Warranty or Additional Liability. While redistributing\n      the Work or Derivative Works thereof, You may choose to offer,\n      and charge a fee for, acceptance of support, warranty, indemnity,\n      or other liability obligations and/or rights consistent with this\n      License. However, in accepting such obligations, You may act only\n      on Your own behalf and on Your sole responsibility, not on behalf\n      of any other Contributor, and only if You agree to indemnify,\n      defend, and hold each Contributor harmless for any liability\n      incurred by, or claims asserted against, such Contributor by reason\n      of your accepting any such warranty or additional liability.\n\n   END OF TERMS AND CONDITIONS\n\n   APPENDIX: How to apply the Apache License to your work.\n\n      To apply the Apache License to your work, attach the following\n      boilerplate notice, with the fields enclosed by brackets \"[]\"\n      replaced with your own identifying information. (Don't include\n      the brackets!)  The text should be enclosed in the appropriate\n      comment syntax for the file format. We also recommend that a\n      file or class name and description of purpose be included on the\n      same \"printed page\" as the copyright notice for easier\n      identification within third-party archives.\n\n   Copyright [yyyy] [name of copyright owner]\n\n   Licensed under the Apache License, Version 2.0 (the \"License\");\n   you may not use this file except in compliance with the License.\n   You may obtain a copy of the License at\n\n       http://www.apache.org/licenses/LICENSE-2.0\n\n   Unless required by applicable law or agreed to in writing, software\n   distributed under the License is distributed on an \"AS IS\" BASIS,\n   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n   See the License for the specific language governing permissions and\n   limitations under the License.\n\n================================================================================\n\n================================================================================\nModule github.com/pkg/errors\nLicensed under the Terms of the BSD-2-Clause License, see also https://github.com/pkg/errors/blob/v0.9.1/LICENSE. \n\nCopyright (c) 2015, Dave Cheney <dave@cheney.net>\nAll rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are met:\n\n* Redistributions of source code must retain the above copyright notice, this\n  list of conditions and the following disclaimer.\n\n* Redistributions in binary form must reproduce the above copyright notice,\n  this list of conditions and the following disclaimer in the documentation\n  and/or other materials provided with the distribution.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\"\nAND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE\nIMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE\nFOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL\nDAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR\nSERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER\nCAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,\nOR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/mattn/go-isatty\nLicensed under the Terms of the MIT License, see also https://github.com/mattn/go-isatty/blob/v0.0.16/LICENSE. \n\nCopyright (c) Yasuhiro MATSUMOTO <mattn.jp@gmail.com>\n\nMIT License (Expat)\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the \"Software\"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/mattn/go-colorable\nLicensed under the Terms of the MIT License, see also https://github.com/mattn/go-colorable/blob/v0.1.13/LICENSE. \n\nThe MIT License (MIT)\n\nCopyright (c) 2016 Yasuhiro Matsumoto\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/julienschmidt/httprouter\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/julienschmidt/httprouter/blob/v1.3.0/LICENSE. \n\nBSD 3-Clause License\n\nCopyright (c) 2013, Julien Schmidt\nAll rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are met:\n\n1. Redistributions of source code must retain the above copyright notice, this\n   list of conditions and the following disclaimer.\n\n2. Redistributions in binary form must reproduce the above copyright notice,\n   this list of conditions and the following disclaimer in the documentation\n   and/or other materials provided with the distribution.\n\n3. Neither the name of the copyright holder nor the names of its\n   contributors may be used to endorse or promote products derived from\n   this software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\"\nAND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE\nIMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE\nFOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL\nDAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR\nSERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER\nCAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,\nOR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/jinzhu/now\nLicensed under the Terms of the MIT License, see also https://github.com/jinzhu/now/blob/v1.1.5/License. \n\nThe MIT License (MIT)\n\nCopyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in\nall copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN\nTHE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/jinzhu/inflection\nLicensed under the Terms of the MIT License, see also https://github.com/jinzhu/inflection/blob/v1.0.0/LICENSE. \n\nThe MIT License (MIT)\n\nCopyright (c) 2015 - Jinzhu\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/jessevdk/go-flags\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/jessevdk/go-flags/blob/v1.5.0/LICENSE. \n\nCopyright (c) 2012 Jesse van den Kieboom. All rights reserved.\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\n     notice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\n     copyright notice, this list of conditions and the following disclaimer\n     in the documentation and/or other materials provided with the\n     distribution.\n   * Neither the name of Google Inc. nor the names of its\n     contributors may be used to endorse or promote products derived from\n     this software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/gosuri/uilive\nLicensed under the Terms of the MIT License, see also https://github.com/gosuri/uilive/blob/v0.0.4/LICENSE. \n\nMIT License\n===========\n\nCopyright (c) 2015, Greg Osuri\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the \"Software\"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/gorilla/websocket\nLicensed under the Terms of the BSD-2-Clause License, see also https://github.com/gorilla/websocket/blob/v1.5.0/LICENSE. \n\nCopyright (c) 2013 The Gorilla WebSocket Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are met:\n\n  Redistributions of source code must retain the above copyright notice, this\n  list of conditions and the following disclaimer.\n\n  Redistributions in binary form must reproduce the above copyright notice,\n  this list of conditions and the following disclaimer in the documentation\n  and/or other materials provided with the distribution.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND\nANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED\nWARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE\nFOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL\nDAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR\nSERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER\nCAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,\nOR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/gorilla/sessions\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/gorilla/sessions/blob/v1.2.1/LICENSE. \n\nCopyright (c) 2012-2018 The Gorilla Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n\t * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n\t * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n\t * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/gorilla/securecookie\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/gorilla/securecookie/blob/v1.1.1/LICENSE. \n\nCopyright (c) 2012 Rodrigo Moraes. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n\t * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n\t * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n\t * Neither the name of Google Inc. nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/gorilla/csrf\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/gorilla/csrf/blob/v1.7.1/LICENSE. \n\nCopyright (c) 2015-2018, The Gorilla Authors. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without modification,\nare permitted provided that the following conditions are met:\n\n1. Redistributions of source code must retain the above copyright notice, this\nlist of conditions and the following disclaimer.\n\n2. Redistributions in binary form must reproduce the above copyright notice,\nthis list of conditions and the following disclaimer in the documentation and/or\nother materials provided with the distribution.\n\n3. Neither the name of the copyright holder nor the names of its contributors\nmay be used to endorse or promote products derived from this software without\nspecific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND\nANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED\nWARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE\nDISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR\nANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES\n(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;\nLOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON\nANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS\nSOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/go-sql-driver/mysql\nLicensed under the Terms of the MPL-2.0 License, see also https://github.com/go-sql-driver/mysql/blob/v1.6.0/LICENSE. \n\nMozilla Public License Version 2.0\n==================================\n\n1. Definitions\n--------------\n\n1.1. \"Contributor\"\n    means each individual or legal entity that creates, contributes to\n    the creation of, or owns Covered Software.\n\n1.2. \"Contributor Version\"\n    means the combination of the Contributions of others (if any) used\n    by a Contributor and that particular Contributor's Contribution.\n\n1.3. \"Contribution\"\n    means Covered Software of a particular Contributor.\n\n1.4. \"Covered Software\"\n    means Source Code Form to which the initial Contributor has attached\n    the notice in Exhibit A, the Executable Form of such Source Code\n    Form, and Modifications of such Source Code Form, in each case\n    including portions thereof.\n\n1.5. \"Incompatible With Secondary Licenses\"\n    means\n\n    (a) that the initial Contributor has attached the notice described\n        in Exhibit B to the Covered Software; or\n\n    (b) that the Covered Software was made available under the terms of\n        version 1.1 or earlier of the License, but not also under the\n        terms of a Secondary License.\n\n1.6. \"Executable Form\"\n    means any form of the work other than Source Code Form.\n\n1.7. \"Larger Work\"\n    means a work that combines Covered Software with other material, in \n    a separate file or files, that is not Covered Software.\n\n1.8. \"License\"\n    means this document.\n\n1.9. \"Licensable\"\n    means having the right to grant, to the maximum extent possible,\n    whether at the time of the initial grant or subsequently, any and\n    all of the rights conveyed by this License.\n\n1.10. \"Modifications\"\n    means any of the following:\n\n    (a) any file in Source Code Form that results from an addition to,\n        deletion from, or modification of the contents of Covered\n        Software; or\n\n    (b) any new file in Source Code Form that contains any Covered\n        Software.\n\n1.11. \"Patent Claims\" of a Contributor\n    means any patent claim(s), including without limitation, method,\n    process, and apparatus claims, in any patent Licensable by such\n    Contributor that would be infringed, but for the grant of the\n    License, by the making, using, selling, offering for sale, having\n    made, import, or transfer of either its Contributions or its\n    Contributor Version.\n\n1.12. \"Secondary License\"\n    means either the GNU General Public License, Version 2.0, the GNU\n    Lesser General Public License, Version 2.1, the GNU Affero General\n    Public License, Version 3.0, or any later versions of those\n    licenses.\n\n1.13. \"Source Code Form\"\n    means the form of the work preferred for making modifications.\n\n1.14. \"You\" (or \"Your\")\n    means an individual or a legal entity exercising rights under this\n    License. For legal entities, \"You\" includes any entity that\n    controls, is controlled by, or is under common control with You. For\n    purposes of this definition, \"control\" means (a) the power, direct\n    or indirect, to cause the direction or management of such entity,\n    whether by contract or otherwise, or (b) ownership of more than\n    fifty percent (50%) of the outstanding shares or beneficial\n    ownership of such entity.\n\n2. License Grants and Conditions\n--------------------------------\n\n2.1. Grants\n\nEach Contributor hereby grants You a world-wide, royalty-free,\nnon-exclusive license:\n\n(a) under intellectual property rights (other than patent or trademark)\n    Licensable by such Contributor to use, reproduce, make available,\n    modify, display, perform, distribute, and otherwise exploit its\n    Contributions, either on an unmodified basis, with Modifications, or\n    as part of a Larger Work; and\n\n(b) under Patent Claims of such Contributor to make, use, sell, offer\n    for sale, have made, import, and otherwise transfer either its\n    Contributions or its Contributor Version.\n\n2.2. Effective Date\n\nThe licenses granted in Section 2.1 with respect to any Contribution\nbecome effective for each Contribution on the date the Contributor first\ndistributes such Contribution.\n\n2.3. Limitations on Grant Scope\n\nThe licenses granted in this Section 2 are the only rights granted under\nthis License. No additional rights or licenses will be implied from the\ndistribution or licensing of Covered Software under this License.\nNotwithstanding Section 2.1(b) above, no patent license is granted by a\nContributor:\n\n(a) for any code that a Contributor has removed from Covered Software;\n    or\n\n(b) for infringements caused by: (i) Your and any other third party's\n    modifications of Covered Software, or (ii) the combination of its\n    Contributions with other software (except as part of its Contributor\n    Version); or\n\n(c) under Patent Claims infringed by Covered Software in the absence of\n    its Contributions.\n\nThis License does not grant any rights in the trademarks, service marks,\nor logos of any Contributor (except as may be necessary to comply with\nthe notice requirements in Section 3.4).\n\n2.4. Subsequent Licenses\n\nNo Contributor makes additional grants as a result of Your choice to\ndistribute the Covered Software under a subsequent version of this\nLicense (see Section 10.2) or under the terms of a Secondary License (if\npermitted under the terms of Section 3.3).\n\n2.5. Representation\n\nEach Contributor represents that the Contributor believes its\nContributions are its original creation(s) or it has sufficient rights\nto grant the rights to its Contributions conveyed by this License.\n\n2.6. Fair Use\n\nThis License is not intended to limit any rights You have under\napplicable copyright doctrines of fair use, fair dealing, or other\nequivalents.\n\n2.7. Conditions\n\nSections 3.1, 3.2, 3.3, and 3.4 are conditions of the licenses granted\nin Section 2.1.\n\n3. Responsibilities\n-------------------\n\n3.1. Distribution of Source Form\n\nAll distribution of Covered Software in Source Code Form, including any\nModifications that You create or to which You contribute, must be under\nthe terms of this License. You must inform recipients that the Source\nCode Form of the Covered Software is governed by the terms of this\nLicense, and how they can obtain a copy of this License. You may not\nattempt to alter or restrict the recipients' rights in the Source Code\nForm.\n\n3.2. Distribution of Executable Form\n\nIf You distribute Covered Software in Executable Form then:\n\n(a) such Covered Software must also be made available in Source Code\n    Form, as described in Section 3.1, and You must inform recipients of\n    the Executable Form how they can obtain a copy of such Source Code\n    Form by reasonable means in a timely manner, at a charge no more\n    than the cost of distribution to the recipient; and\n\n(b) You may distribute such Executable Form under the terms of this\n    License, or sublicense it under different terms, provided that the\n    license for the Executable Form does not attempt to limit or alter\n    the recipients' rights in the Source Code Form under this License.\n\n3.3. Distribution of a Larger Work\n\nYou may create and distribute a Larger Work under terms of Your choice,\nprovided that You also comply with the requirements of this License for\nthe Covered Software. If the Larger Work is a combination of Covered\nSoftware with a work governed by one or more Secondary Licenses, and the\nCovered Software is not Incompatible With Secondary Licenses, this\nLicense permits You to additionally distribute such Covered Software\nunder the terms of such Secondary License(s), so that the recipient of\nthe Larger Work may, at their option, further distribute the Covered\nSoftware under the terms of either this License or such Secondary\nLicense(s).\n\n3.4. Notices\n\nYou may not remove or alter the substance of any license notices\n(including copyright notices, patent notices, disclaimers of warranty,\nor limitations of liability) contained within the Source Code Form of\nthe Covered Software, except that You may alter any license notices to\nthe extent required to remedy known factual inaccuracies.\n\n3.5. Application of Additional Terms\n\nYou may choose to offer, and to charge a fee for, warranty, support,\nindemnity or liability obligations to one or more recipients of Covered\nSoftware. However, You may do so only on Your own behalf, and not on\nbehalf of any Contributor. You must make it absolutely clear that any\nsuch warranty, support, indemnity, or liability obligation is offered by\nYou alone, and You hereby agree to indemnify every Contributor for any\nliability incurred by such Contributor as a result of warranty, support,\nindemnity or liability terms You offer. You may include additional\ndisclaimers of warranty and limitations of liability specific to any\njurisdiction.\n\n4. Inability to Comply Due to Statute or Regulation\n---------------------------------------------------\n\nIf it is impossible for You to comply with any of the terms of this\nLicense with respect to some or all of the Covered Software due to\nstatute, judicial order, or regulation then You must: (a) comply with\nthe terms of this License to the maximum extent possible; and (b)\ndescribe the limitations and the code they affect. Such description must\nbe placed in a text file included with all distributions of the Covered\nSoftware under this License. Except to the extent prohibited by statute\nor regulation, such description must be sufficiently detailed for a\nrecipient of ordinary skill to be able to understand it.\n\n5. Termination\n--------------\n\n5.1. The rights granted under this License will terminate automatically\nif You fail to comply with any of its terms. However, if You become\ncompliant, then the rights granted under this License from a particular\nContributor are reinstated (a) provisionally, unless and until such\nContributor explicitly and finally terminates Your grants, and (b) on an\nongoing basis, if such Contributor fails to notify You of the\nnon-compliance by some reasonable means prior to 60 days after You have\ncome back into compliance. Moreover, Your grants from a particular\nContributor are reinstated on an ongoing basis if such Contributor\nnotifies You of the non-compliance by some reasonable means, this is the\nfirst time You have received notice of non-compliance with this License\nfrom such Contributor, and You become compliant prior to 30 days after\nYour receipt of the notice.\n\n5.2. If You initiate litigation against any entity by asserting a patent\ninfringement claim (excluding declaratory judgment actions,\ncounter-claims, and cross-claims) alleging that a Contributor Version\ndirectly or indirectly infringes any patent, then the rights granted to\nYou by any and all Contributors for the Covered Software under Section\n2.1 of this License shall terminate.\n\n5.3. In the event of termination under Sections 5.1 or 5.2 above, all\nend user license agreements (excluding distributors and resellers) which\nhave been validly granted by You or Your distributors under this License\nprior to termination shall survive termination.\n\n************************************************************************\n*                                                                      *\n*  6. Disclaimer of Warranty                                           *\n*  -------------------------                                           *\n*                                                                      *\n*  Covered Software is provided under this License on an \"as is\"       *\n*  basis, without warranty of any kind, either expressed, implied, or  *\n*  statutory, including, without limitation, warranties that the       *\n*  Covered Software is free of defects, merchantable, fit for a        *\n*  particular purpose or non-infringing. The entire risk as to the     *\n*  quality and performance of the Covered Software is with You.        *\n*  Should any Covered Software prove defective in any respect, You     *\n*  (not any Contributor) assume the cost of any necessary servicing,   *\n*  repair, or correction. This disclaimer of warranty constitutes an   *\n*  essential part of this License. No use of any Covered Software is   *\n*  authorized under this License except under this disclaimer.         *\n*                                                                      *\n************************************************************************\n\n************************************************************************\n*                                                                      *\n*  7. Limitation of Liability                                          *\n*  --------------------------                                          *\n*                                                                      *\n*  Under no circumstances and under no legal theory, whether tort      *\n*  (including negligence), contract, or otherwise, shall any           *\n*  Contributor, or anyone who distributes Covered Software as          *\n*  permitted above, be liable to You for any direct, indirect,         *\n*  special, incidental, or consequential damages of any character      *\n*  including, without limitation, damages for lost profits, loss of    *\n*  goodwill, work stoppage, computer failure or malfunction, or any    *\n*  and all other commercial damages or losses, even if such party      *\n*  shall have been informed of the possibility of such damages. This   *\n*  limitation of liability shall not apply to liability for death or   *\n*  personal injury resulting from such party's negligence to the       *\n*  extent applicable law prohibits such limitation. Some               *\n*  jurisdictions do not allow the exclusion or limitation of           *\n*  incidental or consequential damages, so this exclusion and          *\n*  limitation may not apply to You.                                    *\n*                                                                      *\n************************************************************************\n\n8. Litigation\n-------------\n\nAny litigation relating to this License may be brought only in the\ncourts of a jurisdiction where the defendant maintains its principal\nplace of business and such litigation shall be governed by laws of that\njurisdiction, without reference to its conflict-of-law provisions.\nNothing in this Section shall prevent a party's ability to bring\ncross-claims or counter-claims.\n\n9. Miscellaneous\n----------------\n\nThis License represents the complete agreement concerning the subject\nmatter hereof. If any provision of this License is held to be\nunenforceable, such provision shall be reformed only to the extent\nnecessary to make it enforceable. Any law or regulation which provides\nthat the language of a contract shall be construed against the drafter\nshall not be used to construe this License against a Contributor.\n\n10. Versions of the License\n---------------------------\n\n10.1. New Versions\n\nMozilla Foundation is the license steward. Except as provided in Section\n10.3, no one other than the license steward has the right to modify or\npublish new versions of this License. Each version will be given a\ndistinguishing version number.\n\n10.2. Effect of New Versions\n\nYou may distribute the Covered Software under the terms of the version\nof the License under which You originally received the Covered Software,\nor under the terms of any subsequent version published by the license\nsteward.\n\n10.3. Modified Versions\n\nIf you create software not governed by this License, and you want to\ncreate a new license for such software, you may create and use a\nmodified version of this License if you rename the license and remove\nany references to the name of the license steward (except to note that\nsuch modified license differs from this License).\n\n10.4. Distributing Source Code Form that is Incompatible With Secondary\nLicenses\n\nIf You choose to distribute Source Code Form that is Incompatible With\nSecondary Licenses under the terms of this version of the License, the\nnotice described in Exhibit B of this License must be attached.\n\nExhibit A - Source Code Form License Notice\n-------------------------------------------\n\n  This Source Code Form is subject to the terms of the Mozilla Public\n  License, v. 2.0. If a copy of the MPL was not distributed with this\n  file, You can obtain one at http://mozilla.org/MPL/2.0/.\n\nIf it is not possible or desirable to put the notice in a particular\nfile, then You may include the notice in a location (such as a LICENSE\nfile in a relevant directory) where a recipient would be likely to look\nfor such a notice.\n\nYou may add additional accurate notices of copyright ownership.\n\nExhibit B - \"Incompatible With Secondary Licenses\" Notice\n---------------------------------------------------------\n\n  This Source Code Form is \"Incompatible With Secondary Licenses\", as\n  defined by the Mozilla Public License, v. 2.0.\n\n================================================================================\n\n================================================================================\nModule github.com/gliderlabs/ssh\nLicensed under the Terms of the BSD-3-Clause License, see also https://github.com/gliderlabs/ssh/blob/v0.3.5/LICENSE. \n\nCopyright (c) 2016 Glider Labs. All rights reserved.\n\nRedistribution and use in source and binary forms, with or without\nmodification, are permitted provided that the following conditions are\nmet:\n\n   * Redistributions of source code must retain the above copyright\nnotice, this list of conditions and the following disclaimer.\n   * Redistributions in binary form must reproduce the above\ncopyright notice, this list of conditions and the following disclaimer\nin the documentation and/or other materials provided with the\ndistribution.\n   * Neither the name of Glider Labs nor the names of its\ncontributors may be used to endorse or promote products derived from\nthis software without specific prior written permission.\n\nTHIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS\n\"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT\nLIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR\nA PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT\nOWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,\nSPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT\nLIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,\nDATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY\nTHEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT\n(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE\nOF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.\n\n================================================================================\n\n================================================================================\nModule github.com/feiin/sqlstring\nLicensed under the Terms of the MIT License, see also https://github.com/feiin/sqlstring/blob/v0.3.0/LICENSE. \n\nMIT License\n\nCopyright (c) 2021 solar\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/boombuler/barcode\nLicensed under the Terms of the MIT License, see also https://github.com/boombuler/barcode/blob/6c824513bacc/LICENSE. \n\nThe MIT License (MIT)\n\nCopyright (c) 2014 Florian Sundermann\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/anmitsu/go-shlex\nLicensed under the Terms of the MIT License, see also https://github.com/anmitsu/go-shlex/blob/38f4b401e2be/LICENSE. \n\nCopyright (c) anmitsu <anmitsu.s@gmail.com>\n\nPermission is hereby granted, free of charge, to any person obtaining\na copy of this software and associated documentation files (the\n\"Software\"), to deal in the Software without restriction, including\nwithout limitation the rights to use, copy, modify, merge, publish,\ndistribute, sublicense, and/or sell copies of the Software, and to\npermit persons to whom the Software is furnished to do so, subject to\nthe following conditions:\n\nThe above copyright notice and this permission notice shall be\nincluded in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND,\nEXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF\nMERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND\nNONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE\nLIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION\nOF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION\nWITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/alessio/shellescape\nLicensed under the Terms of the MIT License, see also https://github.com/alessio/shellescape/blob/v1.4.1/LICENSE. \n\nThe MIT License (MIT)\n\nCopyright (c) 2016 Alessio Treglia\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE.\n\n================================================================================\n\n================================================================================\nModule github.com/Showmax/go-fqdn\nLicensed under the Terms of the Apache-2.0 License, see also https://github.com/Showmax/go-fqdn/blob/v1.0.0/LICENSE. \n\nCopyright since 2015 Showmax s.r.o.\n\nLicensed under the Apache License, Version 2.0 (the \"License\");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n    http://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an \"AS IS\" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n\n================================================================================\n\n================================================================================\nModule github.com/FAU-CDI/wdresolve\nLicensed under the Terms of the AGPL-3.0 License, see also https://github.com/FAU-CDI/wdresolve/blob/c9c6779d7c41/LICENSE. \n\n                    GNU AFFERO GENERAL PUBLIC LICENSE\n                       Version 3, 19 November 2007\n\n Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>\n Everyone is permitted to copy and distribute verbatim copies\n of this license document, but changing it is not allowed.\n\n                            Preamble\n\n  The GNU Affero General Public License is a free, copyleft license for\nsoftware and other kinds of works, specifically designed to ensure\ncooperation with the community in the case of network server software.\n\n  The licenses for most software and other practical works are designed\nto take away your freedom to share and change the works.  By contrast,\nour General Public Licenses are intended to guarantee your freedom to\nshare and change all versions of a program--to make sure it remains free\nsoftware for all its users.\n\n  When we speak of free software, we are referring to freedom, not\nprice.  Our General Public Licenses are designed to make sure that you\nhave the freedom to distribute copies of free software (and charge for\nthem if you wish), that you receive source code or can get it if you\nwant it, that you can change the software or use pieces of it in new\nfree programs, and that you know you can do these things.\n\n  Developers that use our General Public Licenses protect your rights\nwith two steps: (1) assert copyright on the software, and (2) offer\nyou this License which gives you legal permission to copy, distribute\nand/or modify the software.\n\n  A secondary benefit of defending all users' freedom is that\nimprovements made in alternate versions of the program, if they\nreceive widespread use, become available for other developers to\nincorporate.  Many developers of free software are heartened and\nencouraged by the resulting cooperation.  However, in the case of\nsoftware used on network servers, this result may fail to come about.\nThe GNU General Public License permits making a modified version and\nletting the public access it on a server without ever releasing its\nsource code to the public.\n\n  The GNU Affero General Public License is designed specifically to\nensure that, in such cases, the modified source code becomes available\nto the community.  It requires the operator of a network server to\nprovide the source code of the modified version running there to the\nusers of that server.  Therefore, public use of a modified version, on\na publicly accessible server, gives the public access to the source\ncode of the modified version.\n\n  An older license, called the Affero General Public License and\npublished by Affero, was designed to accomplish similar goals.  This is\na different license, not a version of the Affero GPL, but Affero has\nreleased a new version of the Affero GPL which permits relicensing under\nthis license.\n\n  The precise terms and conditions for copying, distribution and\nmodification follow.\n\n                       TERMS AND CONDITIONS\n\n  0. Definitions.\n\n  \"This License\" refers to version 3 of the GNU Affero General Public License.\n\n  \"Copyright\" also means copyright-like laws that apply to other kinds of\nworks, such as semiconductor masks.\n\n  \"The Program\" refers to any copyrightable work licensed under this\nLicense.  Each licensee is addressed as \"you\".  \"Licensees\" and\n\"recipients\" may be individuals or organizations.\n\n  To \"modify\" a work means to copy from or adapt all or part of the work\nin a fashion requiring copyright permission, other than the making of an\nexact copy.  The resulting work is called a \"modified version\" of the\nearlier work or a work \"based on\" the earlier work.\n\n  A \"covered work\" means either the unmodified Program or a work based\non the Program.\n\n  To \"propagate\" a work means to do anything with it that, without\npermission, would make you directly or secondarily liable for\ninfringement under applicable copyright law, except executing it on a\ncomputer or modifying a private copy.  Propagation includes copying,\ndistribution (with or without modification), making available to the\npublic, and in some countries other activities as well.\n\n  To \"convey\" a work means any kind of propagation that enables other\nparties to make or receive copies.  Mere interaction with a user through\na computer network, with no transfer of a copy, is not conveying.\n\n  An interactive user interface displays \"Appropriate Legal Notices\"\nto the extent that it includes a convenient and prominently visible\nfeature that (1) displays an appropriate copyright notice, and (2)\ntells the user that there is no warranty for the work (except to the\nextent that warranties are provided), that licensees may convey the\nwork under this License, and how to view a copy of this License.  If\nthe interface presents a list of user commands or options, such as a\nmenu, a prominent item in the list meets this criterion.\n\n  1. Source Code.\n\n  The \"source code\" for a work means the preferred form of the work\nfor making modifications to it.  \"Object code\" means any non-source\nform of a work.\n\n  A \"Standard Interface\" means an interface that either is an official\nstandard defined by a recognized standards body, or, in the case of\ninterfaces specified for a particular programming language, one that\nis widely used among developers working in that language.\n\n  The \"System Libraries\" of an executable work include anything, other\nthan the work as a whole, that (a) is included in the normal form of\npackaging a Major Component, but which is not part of that Major\nComponent, and (b) serves only to enable use of the work with that\nMajor Component, or to implement a Standard Interface for which an\nimplementation is available to the public in source code form.  A\n\"Major Component\", in this context, means a major essential component\n(kernel, window system, and so on) of the specific operating system\n(if any) on which the executable work runs, or a compiler used to\nproduce the work, or an object code interpreter used to run it.\n\n  The \"Corresponding Source\" for a work in object code form means all\nthe source code needed to generate, install, and (for an executable\nwork) run the object code and to modify the work, including scripts to\ncontrol those activities.  However, it does not include the work's\nSystem Libraries, or general-purpose tools or generally available free\nprograms which are used unmodified in performing those activities but\nwhich are not part of the work.  For example, Corresponding Source\nincludes interface definition files associated with source files for\nthe work, and the source code for shared libraries and dynamically\nlinked subprograms that the work is specifically designed to require,\nsuch as by intimate data communication or control flow between those\nsubprograms and other parts of the work.\n\n  The Corresponding Source need not include anything that users\ncan regenerate automatically from other parts of the Corresponding\nSource.\n\n  The Corresponding Source for a work in source code form is that\nsame work.\n\n  2. Basic Permissions.\n\n  All rights granted under this License are granted for the term of\ncopyright on the Program, and are irrevocable provided the stated\nconditions are met.  This License explicitly affirms your unlimited\npermission to run the unmodified Program.  The output from running a\ncovered work is covered by this License only if the output, given its\ncontent, constitutes a covered work.  This License acknowledges your\nrights of fair use or other equivalent, as provided by copyright law.\n\n  You may make, run and propagate covered works that you do not\nconvey, without conditions so long as your license otherwise remains\nin force.  You may convey covered works to others for the sole purpose\nof having them make modifications exclusively for you, or provide you\nwith facilities for running those works, provided that you comply with\nthe terms of this License in conveying all material for which you do\nnot control copyright.  Those thus making or running the covered works\nfor you must do so exclusively on your behalf, under your direction\nand control, on terms that prohibit them from making any copies of\nyour copyrighted material outside their relationship with you.\n\n  Conveying under any other circumstances is permitted solely under\nthe conditions stated below.  Sublicensing is not allowed; section 10\nmakes it unnecessary.\n\n  3. Protecting Users' Legal Rights From Anti-Circumvention Law.\n\n  No covered work shall be deemed part of an effective technological\nmeasure under any applicable law fulfilling obligations under article\n11 of the WIPO copyright treaty adopted on 20 December 1996, or\nsimilar laws prohibiting or restricting circumvention of such\nmeasures.\n\n  When you convey a covered work, you waive any legal power to forbid\ncircumvention of technological measures to the extent such circumvention\nis effected by exercising rights under this License with respect to\nthe covered work, and you disclaim any intention to limit operation or\nmodification of the work as a means of enforcing, against the work's\nusers, your or third parties' legal rights to forbid circumvention of\ntechnological measures.\n\n  4. Conveying Verbatim Copies.\n\n  You may convey verbatim copies of the Program's source code as you\nreceive it, in any medium, provided that you conspicuously and\nappropriately publish on each copy an appropriate copyright notice;\nkeep intact all notices stating that this License and any\nnon-permissive terms added in accord with section 7 apply to the code;\nkeep intact all notices of the absence of any warranty; and give all\nrecipients a copy of this License along with the Program.\n\n  You may charge any price or no price for each copy that you convey,\nand you may offer support or warranty protection for a fee.\n\n  5. Conveying Modified Source Versions.\n\n  You may convey a work based on the Program, or the modifications to\nproduce it from the Program, in the form of source code under the\nterms of section 4, provided that you also meet all of these conditions:\n\n    a) The work must carry prominent notices stating that you modified\n    it, and giving a relevant date.\n\n    b) The work must carry prominent notices stating that it is\n    released under this License and any conditions added under section\n    7.  This requirement modifies the requirement in section 4 to\n    \"keep intact all notices\".\n\n    c) You must license the entire work, as a whole, under this\n    License to anyone who comes into possession of a copy.  This\n    License will therefore apply, along with any applicable section 7\n    additional terms, to the whole of the work, and all its parts,\n    regardless of how they are packaged.  This License gives no\n    permission to license the work in any other way, but it does not\n    invalidate such permission if you have separately received it.\n\n    d) If the work has interactive user interfaces, each must display\n    Appropriate Legal Notices; however, if the Program has interactive\n    interfaces that do not display Appropriate Legal Notices, your\n    work need not make them do so.\n\n  A compilation of a covered work with other separate and independent\nworks, which are not by their nature extensions of the covered work,\nand which are not combined with it such as to form a larger program,\nin or on a volume of a storage or distribution medium, is called an\n\"aggregate\" if the compilation and its resulting copyright are not\nused to limit the access or legal rights of the compilation's users\nbeyond what the individual works permit.  Inclusion of a covered work\nin an aggregate does not cause this License to apply to the other\nparts of the aggregate.\n\n  6. Conveying Non-Source Forms.\n\n  You may convey a covered work in object code form under the terms\nof sections 4 and 5, provided that you also convey the\nmachine-readable Corresponding Source under the terms of this License,\nin one of these ways:\n\n    a) Convey the object code in, or embodied in, a physical product\n    (including a physical distribution medium), accompanied by the\n    Corresponding Source fixed on a durable physical medium\n    customarily used for software interchange.\n\n    b) Convey the object code in, or embodied in, a physical product\n    (including a physical distribution medium), accompanied by a\n    written offer, valid for at least three years and valid for as\n    long as you offer spare parts or customer support for that product\n    model, to give anyone who possesses the object code either (1) a\n    copy of the Corresponding Source for all the software in the\n    product that is covered by this License, on a durable physical\n    medium customarily used for software interchange, for a price no\n    more than your reasonable cost of physically performing this\n    conveying of source, or (2) access to copy the\n    Corresponding Source from a network server at no charge.\n\n    c) Convey individual copies of the object code with a copy of the\n    written offer to provide the Corresponding Source.  This\n    alternative is allowed only occasionally and noncommercially, and\n    only if you received the object code with such an offer, in accord\n    with subsection 6b.\n\n    d) Convey the object code by offering access from a designated\n    place (gratis or for a charge), and offer equivalent access to the\n    Corresponding Source in the same way through the same place at no\n    further charge.  You need not require recipients to copy the\n    Corresponding Source along with the object code.  If the place to\n    copy the object code is a network server, the Corresponding Source\n    may be on a different server (operated by you or a third party)\n    that supports equivalent copying facilities, provided you maintain\n    clear directions next to the object code saying where to find the\n    Corresponding Source.  Regardless of what server hosts the\n    Corresponding Source, you remain obligated to ensure that it is\n    available for as long as needed to satisfy these requirements.\n\n    e) Convey the object code using peer-to-peer transmission, provided\n    you inform other peers where the object code and Corresponding\n    Source of the work are being offered to the general public at no\n    charge under subsection 6d.\n\n  A separable portion of the object code, whose source code is excluded\nfrom the Corresponding Source as a System Library, need not be\nincluded in conveying the object code work.\n\n  A \"User Product\" is either (1) a \"consumer product\", which means any\ntangible personal property which is normally used for personal, family,\nor household purposes, or (2) anything designed or sold for incorporation\ninto a dwelling.  In determining whether a product is a consumer product,\ndoubtful cases shall be resolved in favor of coverage.  For a particular\nproduct received by a particular user, \"normally used\" refers to a\ntypical or common use of that class of product, regardless of the status\nof the particular user or of the way in which the particular user\nactually uses, or expects or is expected to use, the product.  A product\nis a consumer product regardless of whether the product has substantial\ncommercial, industrial or non-consumer uses, unless such uses represent\nthe only significant mode of use of the product.\n\n  \"Installation Information\" for a User Product means any methods,\nprocedures, authorization keys, or other information required to install\nand execute modified versions of a covered work in that User Product from\na modified version of its Corresponding Source.  The information must\nsuffice to ensure that the continued functioning of the modified object\ncode is in no case prevented or interfered with solely because\nmodification has been made.\n\n  If you convey an object code work under this section in, or with, or\nspecifically for use in, a User Product, and the conveying occurs as\npart of a transaction in which the right of possession and use of the\nUser Product is transferred to the recipient in perpetuity or for a\nfixed term (regardless of how the transaction is characterized), the\nCorresponding Source conveyed under this section must be accompanied\nby the Installation Information.  But this requirement does not apply\nif neither you nor any third party retains the ability to install\nmodified object code on the User Product (for example, the work has\nbeen installed in ROM).\n\n  The requirement to provide Installation Information does not include a\nrequirement to continue to provide support service, warranty, or updates\nfor a work that has been modified or installed by the recipient, or for\nthe User Product in which it has been modified or installed.  Access to a\nnetwork may be denied when the modification itself materially and\nadversely affects the operation of the network or violates the rules and\nprotocols for communication across the network.\n\n  Corresponding Source conveyed, and Installation Information provided,\nin accord with this section must be in a format that is publicly\ndocumented (and with an implementation available to the public in\nsource code form), and must require no special password or key for\nunpacking, reading or copying.\n\n  7. Additional Terms.\n\n  \"Additional permissions\" are terms that supplement the terms of this\nLicense by making exceptions from one or more of its conditions.\nAdditional permissions that are applicable to the entire Program shall\nbe treated as though they were included in this License, to the extent\nthat they are valid under applicable law.  If additional permissions\napply only to part of the Program, that part may be used separately\nunder those permissions, but the entire Program remains governed by\nthis License without regard to the additional permissions.\n\n  When you convey a copy of a covered work, you may at your option\nremove any additional permissions from that copy, or from any part of\nit.  (Additional permissions may be written to require their own\nremoval in certain cases when you modify the work.)  You may place\nadditional permissions on material, added by you to a covered work,\nfor which you have or can give appropriate copyright permission.\n\n  Notwithstanding any other provision of this License, for material you\nadd to a covered work, you may (if authorized by the copyright holders of\nthat material) supplement the terms of this License with terms:\n\n    a) Disclaiming warranty or limiting liability differently from the\n    terms of sections 15 and 16 of this License; or\n\n    b) Requiring preservation of specified reasonable legal notices or\n    author attributions in that material or in the Appropriate Legal\n    Notices displayed by works containing it; or\n\n    c) Prohibiting misrepresentation of the origin of that material, or\n    requiring that modified versions of such material be marked in\n    reasonable ways as different from the original version; or\n\n    d) Limiting the use for publicity purposes of names of licensors or\n    authors of the material; or\n\n    e) Declining to grant rights under trademark law for use of some\n    trade names, trademarks, or service marks; or\n\n    f) Requiring indemnification of licensors and authors of that\n    material by anyone who conveys the material (or modified versions of\n    it) with contractual assumptions of liability to the recipient, for\n    any liability that these contractual assumptions directly impose on\n    those licensors and authors.\n\n  All other non-permissive additional terms are considered \"further\nrestrictions\" within the meaning of section 10.  If the Program as you\nreceived it, or any part of it, contains a notice stating that it is\ngoverned by this License along with a term that is a further\nrestriction, you may remove that term.  If a license document contains\na further restriction but permits relicensing or conveying under this\nLicense, you may add to a covered work material governed by the terms\nof that license document, provided that the further restriction does\nnot survive such relicensing or conveying.\n\n  If you add terms to a covered work in accord with this section, you\nmust place, in the relevant source files, a statement of the\nadditional terms that apply to those files, or a notice indicating\nwhere to find the applicable terms.\n\n  Additional terms, permissive or non-permissive, may be stated in the\nform of a separately written license, or stated as exceptions;\nthe above requirements apply either way.\n\n  8. Termination.\n\n  You may not propagate or modify a covered work except as expressly\nprovided under this License.  Any attempt otherwise to propagate or\nmodify it is void, and will automatically terminate your rights under\nthis License (including any patent licenses granted under the third\nparagraph of section 11).\n\n  However, if you cease all violation of this License, then your\nlicense from a particular copyright holder is reinstated (a)\nprovisionally, unless and until the copyright holder explicitly and\nfinally terminates your license, and (b) permanently, if the copyright\nholder fails to notify you of the violation by some reasonable means\nprior to 60 days after the cessation.\n\n  Moreover, your license from a particular copyright holder is\nreinstated permanently if the copyright holder notifies you of the\nviolation by some reasonable means, this is the first time you have\nreceived notice of violation of this License (for any work) from that\ncopyright holder, and you cure the violation prior to 30 days after\nyour receipt of the notice.\n\n  Termination of your rights under this section does not terminate the\nlicenses of parties who have received copies or rights from you under\nthis License.  If your rights have been terminated and not permanently\nreinstated, you do not qualify to receive new licenses for the same\nmaterial under section 10.\n\n  9. Acceptance Not Required for Having Copies.\n\n  You are not required to accept this License in order to receive or\nrun a copy of the Program.  Ancillary propagation of a covered work\noccurring solely as a consequence of using peer-to-peer transmission\nto receive a copy likewise does not require acceptance.  However,\nnothing other than this License grants you permission to propagate or\nmodify any covered work.  These actions infringe copyright if you do\nnot accept this License.  Therefore, by modifying or propagating a\ncovered work, you indicate your acceptance of this License to do so.\n\n  10. Automatic Licensing of Downstream Recipients.\n\n  Each time you convey a covered work, the recipient automatically\nreceives a license from the original licensors, to run, modify and\npropagate that work, subject to this License.  You are not responsible\nfor enforcing compliance by third parties with this License.\n\n  An \"entity transaction\" is a transaction transferring control of an\norganization, or substantially all assets of one, or subdividing an\norganization, or merging organizations.  If propagation of a covered\nwork results from an entity transaction, each party to that\ntransaction who receives a copy of the work also receives whatever\nlicenses to the work the party's predecessor in interest had or could\ngive under the previous paragraph, plus a right to possession of the\nCorresponding Source of the work from the predecessor in interest, if\nthe predecessor has it or can get it with reasonable efforts.\n\n  You may not impose any further restrictions on the exercise of the\nrights granted or affirmed under this License.  For example, you may\nnot impose a license fee, royalty, or other charge for exercise of\nrights granted under this License, and you may not initiate litigation\n(including a cross-claim or counterclaim in a lawsuit) alleging that\nany patent claim is infringed by making, using, selling, offering for\nsale, or importing the Program or any portion of it.\n\n  11. Patents.\n\n  A \"contributor\" is a copyright holder who authorizes use under this\nLicense of the Program or a work on which the Program is based.  The\nwork thus licensed is called the contributor's \"contributor version\".\n\n  A contributor's \"essential patent claims\" are all patent claims\nowned or controlled by the contributor, whether already acquired or\nhereafter acquired, that would be infringed by some manner, permitted\nby this License, of making, using, or selling its contributor version,\nbut do not include claims that would be infringed only as a\nconsequence of further modification of the contributor version.  For\npurposes of this definition, \"control\" includes the right to grant\npatent sublicenses in a manner consistent with the requirements of\nthis License.\n\n  Each contributor grants you a non-exclusive, worldwide, royalty-free\npatent license under the contributor's essential patent claims, to\nmake, use, sell, offer for sale, import and otherwise run, modify and\npropagate the contents of its contributor version.\n\n  In the following three paragraphs, a \"patent license\" is any express\nagreement or commitment, however denominated, not to enforce a patent\n(such as an express permission to practice a patent or covenant not to\nsue for patent infringement).  To \"grant\" such a patent license to a\nparty means to make such an agreement or commitment not to enforce a\npatent against the party.\n\n  If you convey a covered work, knowingly relying on a patent license,\nand the Corresponding Source of the work is not available for anyone\nto copy, free of charge and under the terms of this License, through a\npublicly available network server or other readily accessible means,\nthen you must either (1) cause the Corresponding Source to be so\navailable, or (2) arrange to deprive yourself of the benefit of the\npatent license for this particular work, or (3) arrange, in a manner\nconsistent with the requirements of this License, to extend the patent\nlicense to downstream recipients.  \"Knowingly relying\" means you have\nactual knowledge that, but for the patent license, your conveying the\ncovered work in a country, or your recipient's use of the covered work\nin a country, would infringe one or more identifiable patents in that\ncountry that you have reason to believe are valid.\n\n  If, pursuant to or in connection with a single transaction or\narrangement, you convey, or propagate by procuring conveyance of, a\ncovered work, and grant a patent license to some of the parties\nreceiving the covered work authorizing them to use, propagate, modify\nor convey a specific copy of the covered work, then the patent license\nyou grant is automatically extended to all recipients of the covered\nwork and works based on it.\n\n  A patent license is \"discriminatory\" if it does not include within\nthe scope of its coverage, prohibits the exercise of, or is\nconditioned on the non-exercise of one or more of the rights that are\nspecifically granted under this License.  You may not convey a covered\nwork if you are a party to an arrangement with a third party that is\nin the business of distributing software, under which you make payment\nto the third party based on the extent of your activity of conveying\nthe work, and under which the third party grants, to any of the\nparties who would receive the covered work from you, a discriminatory\npatent license (a) in connection with copies of the covered work\nconveyed by you (or copies made from those copies), or (b) primarily\nfor and in connection with specific products or compilations that\ncontain the covered work, unless you entered into that arrangement,\nor that patent license was granted, prior to 28 March 2007.\n\n  Nothing in this License shall be construed as excluding or limiting\nany implied license or other defenses to infringement that may\notherwise be available to you under applicable patent law.\n\n  12. No Surrender of Others' Freedom.\n\n  If conditions are imposed on you (whether by court order, agreement or\notherwise) that contradict the conditions of this License, they do not\nexcuse you from the conditions of this License.  If you cannot convey a\ncovered work so as to satisfy simultaneously your obligations under this\nLicense and any other pertinent obligations, then as a consequence you may\nnot convey it at all.  For example, if you agree to terms that obligate you\nto collect a royalty for further conveying from those to whom you convey\nthe Program, the only way you could satisfy both those terms and this\nLicense would be to refrain entirely from conveying the Program.\n\n  13. Remote Network Interaction; Use with the GNU General Public License.\n\n  Notwithstanding any other provision of this License, if you modify the\nProgram, your modified version must prominently offer all users\ninteracting with it remotely through a computer network (if your version\nsupports such interaction) an opportunity to receive the Corresponding\nSource of your version by providing access to the Corresponding Source\nfrom a network server at no charge, through some standard or customary\nmeans of facilitating copying of software.  This Corresponding Source\nshall include the Corresponding Source for any work covered by version 3\nof the GNU General Public License that is incorporated pursuant to the\nfollowing paragraph.\n\n  Notwithstanding any other provision of this License, you have\npermission to link or combine any covered work with a work licensed\nunder version 3 of the GNU General Public License into a single\ncombined work, and to convey the resulting work.  The terms of this\nLicense will continue to apply to the part which is the covered work,\nbut the work with which it is combined will remain governed by version\n3 of the GNU General Public License.\n\n  14. Revised Versions of this License.\n\n  The Free Software Foundation may publish revised and/or new versions of\nthe GNU Affero General Public License from time to time.  Such new versions\nwill be similar in spirit to the present version, but may differ in detail to\naddress new problems or concerns.\n\n  Each version is given a distinguishing version number.  If the\nProgram specifies that a certain numbered version of the GNU Affero General\nPublic License \"or any later version\" applies to it, you have the\noption of following the terms and conditions either of that numbered\nversion or of any later version published by the Free Software\nFoundation.  If the Program does not specify a version number of the\nGNU Affero General Public License, you may choose any version ever published\nby the Free Software Foundation.\n\n  If the Program specifies that a proxy can decide which future\nversions of the GNU Affero General Public License can be used, that proxy's\npublic statement of acceptance of a version permanently authorizes you\nto choose that version for the Program.\n\n  Later license versions may give you additional or different\npermissions.  However, no additional obligations are imposed on any\nauthor or copyright holder as a result of your choosing to follow a\nlater version.\n\n  15. Disclaimer of Warranty.\n\n  THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY\nAPPLICABLE LAW.  EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT\nHOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM \"AS IS\" WITHOUT WARRANTY\nOF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,\nTHE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR\nPURPOSE.  THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM\nIS WITH YOU.  SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF\nALL NECESSARY SERVICING, REPAIR OR CORRECTION.\n\n  16. Limitation of Liability.\n\n  IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING\nWILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MODIFIES AND/OR CONVEYS\nTHE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES, INCLUDING ANY\nGENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF THE\nUSE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED TO LOSS OF\nDATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD\nPARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER PROGRAMS),\nEVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF\nSUCH DAMAGES.\n\n  17. Interpretation of Sections 15 and 16.\n\n  If the disclaimer of warranty and limitation of liability provided\nabove cannot be given local legal effect according to their terms,\nreviewing courts shall apply local law that most closely approximates\nan absolute waiver of all civil liability in connection with the\nProgram, unless a warranty or assumption of liability accompanies a\ncopy of the Program in return for a fee.\n\n                     END OF TERMS AND CONDITIONS\n\n            How to Apply These Terms to Your New Programs\n\n  If you develop a new program, and you want it to be of the greatest\npossible use to the public, the best way to achieve this is to make it\nfree software which everyone can redistribute and change under these terms.\n\n  To do so, attach the following notices to the program.  It is safest\nto attach them to the start of each source file to most effectively\nstate the exclusion of warranty; and each file should have at least\nthe \"copyright\" line and a pointer to where the full notice is found.\n\n    <one line to give the program's name and a brief idea of what it does.>\n    Copyright (C) <year>  <name of author>\n\n    This program is free software: you can redistribute it and/or modify\n    it under the terms of the GNU Affero General Public License as published by\n    the Free Software Foundation, either version 3 of the License, or\n    (at your option) any later version.\n\n    This program is distributed in the hope that it will be useful,\n    but WITHOUT ANY WARRANTY; without even the implied warranty of\n    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n    GNU Affero General Public License for more details.\n\n    You should have received a copy of the GNU Affero General Public License\n    along with this program.  If not, see <https://www.gnu.org/licenses/>.\n\nAlso add information on how to contact you by electronic and paper mail.\n\n  If your software can interact with users remotely through a computer\nnetwork, you should also make sure that it provides a way for users to\nget its source.  For example, if your program is a web application, its\ninterface could display a \"Source\" link that leads users to an archive\nof the code.  There are many ways you could offer source, and different\nsolutions will be better for different programs; see section 13 for the\nspecific requirements.\n\n  You should also get your employer (if you work as a programmer) or school,\nif any, to sign a \"copyright disclaimer\" for the program, if necessary.\nFor more information on this, and how to apply and follow the GNU AGPL, see\n<https://www.gnu.org/licenses/>.\n\n================================================================================\n"
}
