sed -i 's/func TestIsExtensionSupported/func TestPipelineIsExtensionSupported/g' pipeline_test.go
sed -i 's/func TestIsLayerSupported/func TestPipelineIsLayerSupported/g' pipeline_test.go
sed -i 's/TransitionImageLayoutSimple/TransitionImageLayout/g' resources_test.go
