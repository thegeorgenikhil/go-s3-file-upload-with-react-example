import {
  Alert,
  AlertIcon,
  Box,
  Button,
  Center,
  Flex,
  Heading,
  Input,
  Link,
  Stack,
  Text,
  useToast,
} from "@chakra-ui/react";
import { useRef, useState } from "react";

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [url, setUrl] = useState<string | null>(null);
  const toast = useToast();
  const fileRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleClick = async () => {
    if (!file) {
      toast({
        title: "No file selected",
        description: "Please select a file to upload",
        status: "error",
        duration: 3000,
      });
      return;
    }
    if (file) {
      try {
        setIsLoading(true);
        const formData = new FormData();
        formData.append("file", file);
        const response = await fetch("http://localhost:8080/upload", {
          method: "POST",
          body: formData,
        });
        const data = await response.json();
        if (data.url) {
          toast({
            title: "File Uploaded Successfully!",
            status: "success",
          });
          setUrl(data.url);
          if (fileRef.current) {
            fileRef.current.value = "";
          }
        }
      } catch (error) {
        console.log(error);
        toast({
          title: "Couldn't upload file to S3",
          status: "error",
          duration: 3000,
        });
      } finally {
        setFile(null);
        setIsLoading(false);
      }
    }
  };

  return (
    <div className="App">
      <Heading p={10} textAlign={"center"}>
        Upload Files to S3 using Go and React
      </Heading>
      <Flex width={"100vw"} alignContent={"center"} justifyContent={"center"}>
        <Center>
          <Box p={5} shadow={"md"}>
            <Stack spacing={3}>
              <Input ref={fileRef} type="file" onChange={handleFileChange} />
              <Button
                isLoading={isLoading}
                onClick={handleClick}
                colorScheme="green"
              >
                Upload File
              </Button>
            </Stack>
          </Box>
        </Center>
      </Flex>
      {url ? (
        <Text p={8} textAlign={"center"}>
          Your file is now available at:{" "}
          <Link textDecoration={"underline"} href={url}>
            {url}
          </Link>
        </Text>
      ) : null}
    </div>
  );
}

export default App;
